package noodle

import (
	"archive/zip"
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/dyweb/gommon/util/fsutil"
	"github.com/pkg/errors"
)

// GenerateEmbed return code without package and the CODE GENERATED header
func GenerateEmbed(root string) ([]byte, error) {
	var (
		err     error
		ignores *fsutil.Ignores
		dirs    = make(map[string]*EmbedDir)
		files   = make(map[string][]*EmbedFile)
		data    []byte
		lastErr error
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
		//if info.Name() == ignoreFileName {
		//	return
		//}
		if info.IsDir() {
			//log.Infof("add %s to path %s", info.Name(), path)
			dirInfo := newEmbedDir(info)
			dirs[join(path, info.Name())] = dirInfo
			dirs[path].Entries = append(dirs[path].Entries, dirInfo.FileInfo)
			return
		}
		// TODO: error group
		if file, err := newEmbedFile(path, info); err != nil {
			log.Warn(err)
			lastErr = err
		} else {
			files[path] = append(files[path], file)
		}
	})
	if lastErr != nil {
		return nil, lastErr
	}
	//log.Infof("dirs (including root) %d", len(dirs))
	//log.Infof("dirs (excluding root) %d", len(files))
	updateDirectoryInfo(dirs, files)
	if data, err = zipFiles(root, files); err != nil {
		return nil, err
	}
	return renderTemplate(root, dirs, data)
}

func readIgnoreFile(root string) (*fsutil.Ignores, error) {
	var err error
	ignores := fsutil.NewIgnores(nil, nil)
	ignoreFile := join(root, ignoreFileName)
	if fsutil.FileExists(ignoreFile) {
		log.Debugf("found ignore file %s", ignoreFile)
		if ignores, err = ReadIgnoreFile(ignoreFile); err != nil {
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
	// TODO: error group
	var lastErr error
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	for path, files := range flatFiles {
		for _, f := range files {
			//log.Infof("write file %s FileSize %d", f.FileName, len(f.Data))
			lastErr = writeZipFile(w, root, path, f)
		}
	}
	if lastErr != nil {
		return nil, lastErr
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

func renderTemplate(root string, dirs map[string]*EmbedDir, data []byte) ([]byte, error) {
	t, err := template.New("noodleembed").Parse(embedTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse embed template")
	}
	// trim root
	trimmedDirs := make(map[string]*EmbedDir)
	for p, d := range dirs {
		trimmedDirs[strings.TrimLeft(p, root)] = d
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, map[string]interface{}{
		"dir":  trimmedDirs,
		"data": data,
	}); err != nil {
		return nil, errors.Wrap(err, "can't execute template")
	}
	//log.Info(buf.String())
	if b, err := format.Source(buf.Bytes()); err != nil {
		return nil, errors.Wrap(err, "can't format go code")
	} else {
		return b, nil
	}
}
