package main

import (
	"bytes"
	"os"
	"net/http"
	"time"
	"fmt"
)

// hand written data before we figure out how to generate them

var mTime = time.Now()

var _ http.FileSystem = (*FileSystem)(nil)

type FileSystem struct {
	files map[string]*File
}

func NewFs(f ... *File) *FileSystem {
	fs := &FileSystem{
		files: make(map[string]*File, len(f)),
	}
	for _, ff := range f {
		// NOTE: "/" is needed
		fs.files["/"+ff.stat.name] = ff
	}
	return fs
}

func (fs *FileSystem) Open(name string) (http.File, error) {
	fmt.Printf("open %s\n", name)
	fmt.Printf("len of files %d\n", len(fs.files))
	if f, exists := fs.files[name]; exists {
		return f, nil
	}
	return nil, os.ErrNotExist
}

// A FileInfo describes a file and is returned by Stat and Lstat.
//type FileInfo interface {
//	Name() string       // base name of the file
//	Size() int64        // length in bytes for regular files; system-dependent for others
//	Mode() FileMode     // file mode bits
//	ModTime() time.Time // modification time
//	IsDir() bool        // abbreviation for Mode().IsDir()
//	Sys() interface{}   // underlying data source (can return nil)
//}

var _ os.FileInfo = (*Info)(nil)

type Info struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (i *Info) Name() string {
	return i.name
}

func (i *Info) Size() int64 {
	return i.size
}

func (i *Info) Mode() os.FileMode {
	return i.mode
}

func (i *Info) ModTime() time.Time {
	return i.modTime
}

func (i *Info) IsDir() bool {
	return i.isDir
}

func (i *Info) Sys() interface{} {
	return nil
}

//type File interface {
//	io.Closer
//	io.Reader
//	io.Seeker
//	Readdir(count int) ([]os.FileInfo, error)
//	Stat() (os.FileInfo, error)
//}

var _ http.File = (*File)(nil)

type File struct {
	stat     Info
	content  []byte
	reader   *bytes.Reader
	allowDir bool
}

func NewFile(name string, b []byte, allowDir bool) *File {
	return &File{
		content:  b,
		reader:   bytes.NewReader(b),
		allowDir: allowDir,
		stat: Info{
			name:    name,
			size:    int64(len(b)),
			mode:    0666,
			modTime: mTime,
			isDir:   false,
		},
	}
}

func (f *File) Read(p []byte) (int, error) {
	return f.reader.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return &f.stat, nil
}

// TODO: what is the count for?
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	// TODO: return dir when it is allowed
	return make([]os.FileInfo, 0), nil
}

func (f *File) Close() error {
	return nil
}
