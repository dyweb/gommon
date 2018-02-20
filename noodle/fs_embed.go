package noodle

import (
	"archive/zip"
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
	info FileInfo
	data []byte
}

// TODO: maybe we can just keep the bytes inside FileInfo, instead of write it to zip along the way ...
type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	//entries []FileInfo
}

func GenerateEmbed(root string) error {
	var (
		err     error
		ignores *fsutil.Ignores
		files   = make(map[string][]*embedFile)
		lastErr error
	)
	if ignores, err = readIgnoreFile(root); err != nil {
		return err
	}
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		// TODO: aggregate the error, errors group?
		//log.Info(path)
		if file, err := newEmbedFile(path, info); err != nil {
			log.Warn(err)
			lastErr = err
		} else {
			files[path] = append(files[path], file)
		}
	})
	log.Info(len(files))
	updateDirectoryInfo(files)
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

func newEmbedFile(path string, info os.FileInfo) (*embedFile, error) {
	var (
		b   []byte
		err error
	)
	if !info.IsDir() {
		b, err = ioutil.ReadFile(join(path, info.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "can't read file from disk")
		}
	}
	return &embedFile{
		info: FileInfo{
			name:  info.Name(),
			size:  info.Size(),
			mode:  info.Mode(),
			isDir: info.IsDir(),
		},
		data: b,
	}, nil
}

func updateDirectoryInfo(flatFiles map[string][]*embedFile) {
	for path, files := range flatFiles {
		for _, f := range files {
			//flatFiles[path]
			// dir size is 4096 4KB ...
			log.Info(f.info.name, " ", f.info.size, f.info.isDir)
		}
		log.Info(path, " ", len(files))
	}
}

func writeZipFile(w *zip.Writer, root string, path string, info os.FileInfo) error {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.Wrap(err, "can't create file header")
	}
	header.Method = zip.Deflate
	header.Name = strings.TrimLeft(join(path, info.Name()), root)
	f, err := w.CreateHeader(header)
	if err != nil {
		return errors.Wrap(err, "can't add file to zip")
	}
	b, err := ioutil.ReadFile(join(path, info.Name()))
	if err != nil {
		return errors.Wrap(err, "can't read file from disk")
	}
	if _, err := f.Write(b); err != nil {
		return errors.Wrap(err, "can't write zip file content")
	}
	return nil
}
