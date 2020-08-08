package fsutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dyweb/gommon/errors"
)

const (
	DefaultFilePerm = 0664
	DefaultDirPerm  = 0775
)

func IsGoFile(info os.FileInfo) bool {
	name := info.Name()
	// not a folder & not hidden & .go
	return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

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

// CreateFileAndPath creates the folder if it does not exists and create a new file using os.Create.
func CreateFileAndPath(path, file string) (*os.File, error) {
	if err := MkdirIfNotExists(path); err != nil {
		return nil, err
	}
	return os.Create(filepath.Join(path, file))
}

// WriteTempFile creates a temporary file and writes data to it.
// Dir can be empty string. It returns path of the created file.
// NOTE: It is based on goimports command.
func WriteTempFile(dir, prefix string, data []byte) (string, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	_, err = f.Write(data)
	if errClose := f.Close(); err == nil {
		err = errClose
	}
	if err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// WriteTempFiles creates multiple files under same dir w/ same prefix.
// It stops if there are any error and always returns created files.
func WriteTempFiles(dir, prefix string, dataList ...[]byte) ([]string, error) {
	var names []string
	for _, data := range dataList {
		n, err := WriteTempFile(dir, prefix, data)
		if err != nil {
			return names, err
		}
		names = append(names, n)
	}
	return names, nil
}

// RemoveFiles removes multiple files. It keeps track of error for each removal.
// But it does NOT stop when there is error.
// Returned non nil error must be a dyweb/gommon/errors.MultiErr.
func RemoveFiles(names []string) error {
	merr := errors.NewMultiErr()
	for _, name := range names {
		merr.Append(os.Remove(name))
	}
	return merr.ErrorOrNil()
}

// Diff compares two files by shelling out to system diff binary.
// NOTE: It is based on goimports command.
// TODO: there are pure go diff package, and there is also diff w/ syntax highlight written in rust
// TODO: allow force color output, default is auto and I guess it detects tty
func Diff(p1 string, p2 string) ([]byte, error) {
	b, err := exec.Command("diff", "-u", p1, p2).CombinedOutput()
	// NOTE: diff returns 1 when there are diff, so we ignore the error as long as there is valid output.
	if len(b) != 0 {
		return b, nil
	}
	return b, errors.Wrapf(err, "error shell out to diff -u %s %s", p1, p2)
}
