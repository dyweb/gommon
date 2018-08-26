package noodle

import (
	"archive/zip"
	"io/ioutil"
	"path/filepath"

	"github.com/dyweb/gommon/errors"
)

func join(elem ...string) string {
	return filepath.Join(elem...)
}

func unzip(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, errors.Wrap(err, "can't open file inside zip")
	}
	if b, err := ioutil.ReadAll(r); err != nil {
		return nil, errors.Wrap(err, "can't read file content")
	} else {
		r.Close()
		return b, nil
	}
}
