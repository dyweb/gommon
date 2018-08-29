package noodle

import (
	"net/http"
	"os"
)

var (
	_ http.FileSystem = (*LocalFs)(nil)
	_ Bowel           = (*LocalFs)(nil)
)

type LocalFs struct {
	root    string
	dir     http.Dir
	listDir bool
}

// NewLocal returns a http.FileSystem using local directory, both list directory and index.html are disabled
// it is a wrapper around http.Dir
func NewLocal(root string) *LocalFs {
	return &LocalFs{root: root, dir: http.Dir(root)}
}

// NewLocalUnsafe allows list directory and use http.Dir directly
func NewLocalUnsafe(root string) *LocalFs {
	fs := NewLocal(root)
	fs.listDir = true
	return fs
}

// Open implements http.FileSystem interface, if list directory is not allowed,
// it will return os.ErrNotExist when it detects the file is a directory
func (fs *LocalFs) Open(name string) (http.File, error) {
	// NOTE: http.Dir has extra error handling, https://github.com/golang/go/issues/18984
	// some operation are mapped to 404 instead of 500
	f, err := fs.dir.Open(name)
	if err != nil {
		return f, err
	}
	// unsafe mode, return without checking if the file is a dir
	if fs.listDir {
		return f, nil
	}
	info, err := f.Stat()
	if err != nil {
		return f, err
	}
	// NOTE: since we disable list directory this way, index.html is no longer supported
	if info.IsDir() {
		return nil, os.ErrNotExist
	}
	// safe mode && not dir
	return f, err
}
