package noodle

import (
	"net/http"
)

var _ http.FileSystem = (*LocalFs)(nil)

type LocalFs struct {
	root string
	dir  http.Dir
}

func NewLocal(root string) *LocalFs {
	return &LocalFs{root: root, dir: http.Dir(root)}
}

func (fs *LocalFs) Open(name string) (http.File, error) {
	//return os.Open(filepath.Join(fs.root, name))
	// NOTE: http.Dir has extra error handling, https://github.com/golang/go/issues/18984
	// some operation are mapped to 404 instead of 500
	// TODO: but it supports list dir and can't be disabled ...
	return fs.dir.Open(name)
}
