package noodle

import (
	"archive/zip"
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/dyweb/gommon/util/genutil"
)

const generator = "gommon/noodle"

type EmbedConfig struct {
	Src     string `json:"src" yaml:"src"`
	Dst     string `json:"dst" yaml:"dst"`
	Name    string `json:"name" yaml:"name"`
	Package string `json:"package" yaml:"package"`
}

// GenerateEmbedFile generates contents from cfg.Src and save it to a single go file in cfg.Dst.
// It's a wrapper around GenerateEmbedBytes
func GenerateEmbedFile(cfg EmbedConfig) error {
	if cfg.Dst == "" {
		return errors.New("empty dst in config")
	}
	b, err := GenerateEmbedBytes([]EmbedConfig{cfg})
	if err != nil {
		return nil
	}
	return fsutil.WriteFile(cfg.Dst, b)
}

// GenerateEmbedBytes return a single formatted go file as bytes from multiple source directories.
// Use GenerateEmbedFile if you just have one source directory and want to write to file directly.
// it's header + package + []GenerateEmbedPartial
func GenerateEmbedBytes(cfgs []EmbedConfig) ([]byte, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("config list is empty")
	}
	pkg := cfgs[0].Package
	var srcs []string
	for _, cfg := range cfgs {
		if cfg.Src == "" {
			return nil, errors.New("empty src in config")
		}
		if pkg != cfg.Package {
			return nil, errors.Errorf("different package name among configs, %s != %s", pkg, cfg.Package)
		}
		srcs = append(srcs, cfg.Src)
	}
	var buf bytes.Buffer
	buf.WriteString(genutil.Header(generator, strings.Join(srcs, ",")))
	buf.WriteString("package " + pkg)
	buf.WriteString(`
import (
	"time"
	
	"github.com/dyweb/gommon/noodle"
)
`)
	for _, cfg := range cfgs {
		if b, err := GenerateEmbedPartial(cfg); err != nil {
			return nil, err
		} else {
			buf.Write(b)
		}
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "can't format go code")
	}
	return b, nil
}

// GenerateEmbedPartial return code WITHOUT header and import, use GenerateEmbedBytes if you want a full go file
func GenerateEmbedPartial(cfg EmbedConfig) ([]byte, error) {
	root := cfg.Src
	var (
		err     error
		ignores *fsutil.Ignores
		dirs    = make(map[string]*EmbedDir)
		files   = make(map[string][]*EmbedFile)
		data    []byte
		merr    = errors.NewMultiErr()
	)
	if rootStat, err := os.Stat(root); err != nil {
		return nil, errors.Wrap(err, "can't get stat of root folder")
	} else {
		//log.Infof("root %s rootStat name %s", root, rootStat.Name())
		dirs[root] = newEmbedDir(rootStat)
	}
	if ignores, err = readIgnoreFile(root); err != nil {
		return nil, err
	}
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		//log.Info(path)
		// TODO: allow config if ignore file should be included
		//if info.Name() == DefaultIgnoreFileName {
		//	return
		//}
		if info.IsDir() {
			//log.Infof("add %s to path %s", info.Name(), path)
			dirInfo := newEmbedDir(info)
			dirs[join(path, info.Name())] = dirInfo
			dirs[path].Entries = append(dirs[path].Entries, dirInfo.FileInfo)
			return
		}
		if file, err := newEmbedFile(path, info); err != nil {
			merr.Append(err)
		} else {
			files[path] = append(files[path], file)
		}
	})
	if merr.HasError() {
		return nil, merr
	}
	//log.Infof("dirs (including root) %d", len(dirs))
	//log.Infof("dirs (excluding root) %d", len(files))
	updateDirectoryInfo(dirs, files)
	if data, err = zipFiles(root, files); err != nil {
		return nil, err
	}
	return renderTemplate(cfg, dirs, data)
}

func readIgnoreFile(root string) (*fsutil.Ignores, error) {
	var err error
	ignores := fsutil.NewIgnores(nil, nil)
	ignoreFile := join(root, DefaultIgnoreFileName)
	if fsutil.FileExists(ignoreFile) {
		log.Debugf("found ignore file %s", ignoreFile)
		if ignores, err = fsutil.ReadIgnoreFile(ignoreFile); err != nil {
			return ignores, err
		}
		// set common prefix so ignore path would work
		ignores.SetPathPrefix(root)
		log.Debugf("ignore patterns %v", ignores.Patterns())
	}
	return ignores, nil
}

func newEmbedDir(info os.FileInfo) *EmbedDir {
	return &EmbedDir{
		FileInfo: FileInfo{
			FileName:    info.Name(),
			FileSize:    info.Size(),
			FileMode:    info.Mode(),
			FileModTime: info.ModTime(),
			FileIsDir:   info.IsDir(),
		},
	}
}

func newEmbedFile(path string, info os.FileInfo) (*EmbedFile, error) {
	if b, err := ioutil.ReadFile(join(path, info.Name())); err != nil {
		return nil, errors.Wrap(err, "can't read file from disk")
	} else {
		return &EmbedFile{
			FileInfo: FileInfo{
				FileName:    info.Name(),
				FileSize:    info.Size(),
				FileMode:    info.Mode(),
				FileModTime: info.ModTime(),
				FileIsDir:   info.IsDir(),
			},
			Data: b,
		}, nil
	}
}

func updateDirectoryInfo(dirs map[string]*EmbedDir, flatFiles map[string][]*EmbedFile) {
	for path, files := range flatFiles {
		//log.Infof("path %s files %d", path, len(files))
		for _, f := range files {
			//log.Infof("add %s to path %s", f.Name(), path)
			dirs[path].Entries = append(dirs[path].Entries, f.FileInfo)
		}
	}
}

func zipFiles(root string, flatFiles map[string][]*EmbedFile) ([]byte, error) {
	merr := errors.NewMultiErr()
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	// sort dir to make the output stable https://github.com/dyweb/gommon/issues/52
	var paths []string
	for path := range flatFiles {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		for _, f := range flatFiles[path] {
			//log.Infof("write file %s FileSize %d", f.FileName, len(f.Data))
			merr.Append(writeZipFile(w, root, path, f))
		}
	}
	if merr.HasError() {
		return nil, merr
	}
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(err, "can't close zip writer")
	}
	return buf.Bytes(), nil
}

func writeZipFile(w *zip.Writer, root string, path string, file *EmbedFile) error {
	header, err := zip.FileInfoHeader(&file.FileInfo)
	if err != nil {
		return errors.Wrap(err, "can't create file header")
	}
	header.Method = zip.Deflate
	// trim root
	header.Name = join(strings.TrimLeft(path, root), file.Name())
	f, err := w.CreateHeader(header)
	if err != nil {
		return errors.Wrap(err, "can't add file to zip")
	}
	if _, err := f.Write(file.Data); err != nil {
		return errors.Wrap(err, "can't write zip file content")
	}
	return nil
}

func renderTemplate(cfg EmbedConfig, dirs map[string]*EmbedDir, data []byte) ([]byte, error) {
	root := cfg.Src
	if cfg.Name == "" {
		log.Debug("export name not set, using default")
		cfg.Name = DefaultName
	}
	t, err := template.New("noodleembed").Parse(embedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse embed template")
	}
	// trim root
	trimmedDirs := make(map[string]*EmbedDir)
	for p, d := range dirs {
		// NOTE: use TrimPrefix instead TrimLeft
		trimmedDirs[strings.TrimPrefix(p, root)] = d
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, map[string]interface{}{
		"dir":  trimmedDirs,
		"data": data,
		"name": cfg.Name,
		"src":  cfg.Src,
	}); err != nil {
		return nil, errors.Wrap(err, "can't execute template")
	}
	//log.Info(buf.String())
	return buf.Bytes(), nil
}
