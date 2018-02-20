package noodle

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/dyweb/gommon/util/fsutil"
	"github.com/pkg/errors"
)

// TODO
// walk the folder, keep track of folders
type embedFile struct {
	FileInfo
	data []byte
}

type embedDir struct {
	FileInfo
	entries []FileInfo
}

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func GenerateEmbed(root string) error {
	var (
		err     error
		ignores *fsutil.Ignores
		dirs    = make(map[string]*embedDir)
		files   = make(map[string][]*embedFile)
		lastErr error
	)
	if rootStat, err := os.Stat(root); err != nil {
		return errors.Wrap(err, "can't get stat of root folder")
	} else {
		log.Infof("root %s rootStat name %s", root, rootStat.Name())
		dirs[root] = newEmbedDir(rootStat)
	}
	if ignores, err = readIgnoreFile(root); err != nil {
		return err
	}
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		//log.Info(path)
		// TODO: allow config if ignore file should be included
		//if info.Name() == ignoreFileName {
		//	return
		//}
		if info.IsDir() {
			dirInfo := newEmbedDir(info)
			dirs[join(path, info.Name())] = dirInfo
			dirs[path].entries = append(dirs[path].entries, dirInfo.FileInfo)
			return
		}
		if file, err := newEmbedFile(path, info); err != nil {
			// TODO: aggregate the error, errors group?
			log.Warn(err)
			lastErr = err
		} else {
			files[path] = append(files[path], file)
		}
	})
	log.Infof("total dirs (including root) %d", len(dirs))
	log.Infof("dirs %d", len(files))
	updateDirectoryInfo(dirs, files)
	err = writeZipFiles("t.zip", root, files)
	log.Warn(err)
	return lastErr
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

func newEmbedDir(info os.FileInfo) *embedDir {
	return &embedDir{
		FileInfo: FileInfo{
			name:  info.Name(),
			size:  info.Size(),
			mode:  info.Mode(),
			isDir: info.IsDir(),
		},
	}
}

func newEmbedFile(path string, info os.FileInfo) (*embedFile, error) {
	if b, err := ioutil.ReadFile(join(path, info.Name())); err != nil {
		return nil, errors.Wrap(err, "can't read file from disk")
	} else {
		return &embedFile{
			FileInfo: FileInfo{
				name:  info.Name(),
				size:  info.Size(),
				mode:  info.Mode(),
				isDir: info.IsDir(),
			},
			data: b,
		}, nil
	}
}

func updateDirectoryInfo(dirs map[string]*embedDir, flatFiles map[string][]*embedFile) {
	//folders := make(map[string][]FileInfo, len(flatFiles) + 1)
	for path, files := range flatFiles {
		log.Infof("path %s files %d", path, len(files))
		for _, f := range files {
			dirs[path].entries = append(dirs[path].entries, f.FileInfo)
		}
	}
}

func writeZipFiles(dst string, root string, flatFiles map[string][]*embedFile) error {
	var lastErr error
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	for path, files := range flatFiles {
		for _, f := range files {
			log.Infof("write file %s size %d", f.name, len(f.data))
			lastErr = writeZipFile(w, root, path, f)
		}
	}

	if lastErr != nil {
		return lastErr
	}
	if err := w.Close(); err != nil {
		return errors.Wrap(err, "can't close zip writer")
	}
	if err := ioutil.WriteFile(dst, buf.Bytes(), 0666); err != nil {
		return errors.Wrap(err, "can't write zip file")
	}
	return nil
}

func writeZipFile(w *zip.Writer, root string, path string, file *embedFile) error {
	header, err := zip.FileInfoHeader(&file.FileInfo)
	if err != nil {
		return errors.Wrap(err, "can't create file header")
	}
	header.Method = zip.Deflate
	header.Name = strings.TrimLeft(join(path, file.FileInfo.Name()), root)
	f, err := w.CreateHeader(header)
	if err != nil {
		return errors.Wrap(err, "can't add file to zip")
	}
	if _, err := f.Write(file.data); err != nil {
		return errors.Wrap(err, "can't write zip file content")
	}
	return nil
}

func (i *FileInfo) Name() string {
	return i.name
}

func (i *FileInfo) Size() int64 {
	return i.size
}

func (i *FileInfo) Mode() os.FileMode {
	return i.mode
}

func (i *FileInfo) ModTime() time.Time {
	return i.modTime
}

func (i *FileInfo) IsDir() bool {
	return i.isDir
}

func (i *FileInfo) Sys() interface{} {
	return nil
}
