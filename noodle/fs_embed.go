package noodle

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dyweb/gommon/util/fsutil"
	"github.com/pkg/errors"
)

// TODO
// walk the folder, keep track of folders
type embedData struct {
	info os.FileInfo
	data []byte
}

func GenerateEmbed(root string) error {
	var (
		err     error
		ignores *fsutil.Ignores
	)
	ignoreFile := join(root, ignoreFileName)
	if fsutil.FileExists(ignoreFile) {
		log.Debugf("found ignore file %s", ignoreFile)
		if ignores, err = ReadIgnoreFile(ignoreFile); err != nil {
			return err
		}
		// set common prefix so ignore path would work
		ignores.SetPathPrefix(root)
		log.Debug(ignores.Patterns())
	}
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		// TODO: register meta, we need this when read dir, which is not supported by statik
		if info.IsDir() {
			return
		}
		log.Info(join(path, info.Name()))
		// TODO: read file and render template, could put it into a single go file with a large byte slice
		// TODO: aggregate the error, errors group?
		if err := writeZipFile(w, root, path, info); err != nil {
			log.Warnf("can't create %v", err)
		}
	})
	w.Close()
	log.Info(buf.Len())
	if err := ioutil.WriteFile("t.zip", buf.Bytes(), 0666); err != nil {
		log.Warnf("can't write zip %v", err)
	}
	return nil
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
