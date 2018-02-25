package fsutil

import (
	"io/ioutil"

	"github.com/dyweb/gommon/errors"
)

// WriteFile use 0666 as permission and wrap standard error
func WriteFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, 0664); err != nil {
		return errors.Wrap(err, "can't write file")
	}
	return nil
}
