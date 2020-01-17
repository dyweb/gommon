package fsutil

import (
	"io/ioutil"
	"os"

	"github.com/dyweb/gommon/errors"
)

const (
	DefaultFilePerm = 0664
	DefaultDirPerm  = 0775
)

// WriteFile use 0664 as permission and wrap standard error
func WriteFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, DefaultFilePerm); err != nil {
		return errors.Wrap(err, "can't write file")
	}
	return nil
}

// MkdirIfNotExists check if the directory already exists before calling os.MkdirAll with perm 0775
func MkdirIfNotExists(path string) error {
	i, err := os.Stat(path)
	// file exists
	if err == nil {
		if i.IsDir() {
			return nil
		}
		return errors.New("path to create dir is a file already: " + path)
	}
	return os.MkdirAll(path, DefaultDirPerm)
}
