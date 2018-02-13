package noodle

import (
	"net/http"
	"os"
)

var _ http.FileSystem = (*LocalFs)(nil)

type LocalFs struct {
	root    string
	dir     http.Dir
	listDir bool
}

type disabledDir struct {
	http.File
}

func (d *disabledDir) Readdir(count int) ([]os.FileInfo, error) {
	// TODO: the draw back is the result page is empty instead of a standard 404 page, which tells user this is a folder ...
	return nil, nil
}

// NewLocal returns a http.FileSystem using local directory, list directory is disabled, it's just a wrap around http.Dir
func NewLocal(root string) *LocalFs {
	return &LocalFs{root: root, dir: http.Dir(root)}
}

// NewLocalUnsafe allows list directory and use http.Dir directly
func NewLocalUnsafe(root string) *LocalFs {
	fs := NewLocal(root)
	fs.listDir = true
	return fs
}

func (fs *LocalFs) Open(name string) (http.File, error) {
	//return os.Open(filepath.Join(fs.root, name))
	// NOTE: http.Dir has extra error handling, https://github.com/golang/go/issues/18984
	// some operation are mapped to 404 instead of 500
	f, err := fs.dir.Open(name)
	if err != nil {
		return f, err
	}
	// NOTE: http.FileServer supports list dir and can't be disabled, this is hack to ban it from list dir
	// the draw back is the result is not a standard 404 error page, but a blank 200 page ...
	if !fs.listDir {
		f = &disabledDir{f}
	}
	return f, err
}
