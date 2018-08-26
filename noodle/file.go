package noodle

import (
	"os"
	"time"
)

var _ os.FileInfo = (*FileInfo)(nil)

// FileInfo is a concrete struct that implements os.FileInfo interface
// Its fields have the awkward File* prefix to avoid conflict with os.FileInfo interface
type FileInfo struct {
	FileName    string
	FileSize    int64
	FileMode    os.FileMode
	FileModTime time.Time
	FileIsDir   bool
}

func (i *FileInfo) Name() string {
	return i.FileName
}

func (i *FileInfo) Size() int64 {
	return i.FileSize
}

func (i *FileInfo) Mode() os.FileMode {
	return i.FileMode
}

func (i *FileInfo) ModTime() time.Time {
	return i.FileModTime
}

func (i *FileInfo) IsDir() bool {
	return i.FileIsDir
}

func (i *FileInfo) Sys() interface{} {
	return nil
}
