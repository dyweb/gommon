package noodle

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dyweb/gommon/errors"
)

var registeredBoxes map[string]EmbedBox

var _ http.FileSystem = (*EmbedBox)(nil)
var _ http.File = (*EmbedDir)(nil)

type EmbedBox struct {
	Dirs  map[string]EmbedDir
	Data  []byte
	files map[string]EmbedFile
}

type EmbedFile struct {
	FileInfo
	Data   []byte
	reader *bytes.Reader
}

type EmbedDir struct {
	FileInfo
	Entries []FileInfo
}

func (d *EmbedDir) Read(p []byte) (int, error) {
	return 0, os.ErrNotExist
}

func (d *EmbedDir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrNotExist
}

func (d *EmbedDir) Stat() (os.FileInfo, error) {
	return &d.FileInfo, nil
}

func (d *EmbedDir) Close() error {
	return nil
}

func (d *EmbedDir) Readdir(count int) ([]os.FileInfo, error) {
	// TODO: disable list dir
	log.Infof("readdir %d", count)
	files := make([]os.FileInfo, 0, len(d.Entries))
	// FIXED: learn this the hard way .. https://github.com/dyweb/gommon/issues/50
	// the element is range syntax is created when loop start and it is reused between iterations, thus same pointer
	// https://tam7t.com/golang-range-and-pointers/
	//for _, f := range d.Entries {
	//	log.Infof("file %s", f.Name())
	//	files = append(files, &f)
	//}
	for i := range d.Entries {
		files = append(files, &d.Entries[i])
	}
	return files, nil
}

func (b *EmbedBox) Open(name string) (http.File, error) {
	// check dir first
	log.Infof("open %s", name)
	// trim /
	name = name[1:]
	if d, exists := b.Dirs[name]; exists {
		log.Infof("%s entries %d", name, len(d.Entries))
		return &d, nil
	}
	// check file
	if f, exists := b.files[name]; exists {
		return &f, nil
	}
	return nil, os.ErrNotExist
}

func (b *EmbedBox) ExtractFiles() error {
	r, err := zip.NewReader(bytes.NewReader(b.Data), int64(len(b.Data)))
	if err != nil {
		errors.Wrap(err, "can't read zipped data")
	}
	b.files = make(map[string]EmbedFile)
	for _, f := range r.File {
		if bs, err := unzip(f); err != nil {
			return err
		} else {
			b.files[f.Name] = EmbedFile{
				FileInfo: *NewFileInfo(f.FileInfo()),
				Data:     bs,
			}
		}
	}
	return nil
}

func (f *EmbedFile) Read(p []byte) (int, error) {
	if f.reader == nil {
		f.reader = bytes.NewReader(f.Data)
	}
	return f.reader.Read(p)
}

func (f *EmbedFile) Seek(offset int64, whence int) (int64, error) {
	if f.reader == nil {
		f.reader = bytes.NewReader(f.Data)
	}
	return f.reader.Seek(offset, whence)
}

func (f *EmbedFile) Stat() (os.FileInfo, error) {
	return &f.FileInfo, nil
}

func (f *EmbedFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrNotExist
}

func (f *EmbedFile) Close() error {
	// TODO: I think we can do nothing?
	//if f.reader != nil {
	//	f.reader = nil
	//}
	return nil
}

func RegisterEmbedBox(name string, box EmbedBox) {
	log.Debugf("register embed box %s", name)
	if _, exists := registeredBoxes[name]; exists {
		log.Warnf("box %s already exists, overwrite it now", name)
	}
	registeredBoxes[name] = box
}

func GetEmbedBox(name string) (EmbedBox, error) {
	if box, exists := registeredBoxes[name]; exists {
		return box, nil
	} else {
		return box, errors.Errorf("box %s does not exist", name)
	}
}

func NewFileInfo(info os.FileInfo) *FileInfo {
	return &FileInfo{
		FileName:    info.Name(),
		FileSize:    info.Size(),
		FileMode:    info.Mode(),
		FileModTime: info.ModTime(),
		FileIsDir:   info.IsDir(),
	}
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

func init() {
	registeredBoxes = make(map[string]EmbedBox)
}
