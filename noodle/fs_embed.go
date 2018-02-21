package noodle

import (
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

var registeredBoxes map[string]EmbedBox

var _ http.FileSystem = (*EmbedBox)(nil)
var _ http.File = (*EmbedDir)(nil)

type EmbedBox struct {
	Dirs map[string]EmbedDir
	Data []byte
}

type EmbedFile struct {
	FileInfo
	Data []byte
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
	files := make([]os.FileInfo, len(d.Entries))
	for _, f := range d.Entries {
		files = append(files, &f)
	}
	return files, nil
}

func (b *EmbedBox) Open(name string) (http.File, error) {
	// check dir first
	log.Infof("open %s", name)
	name = name[1:]
	if d, exists := b.Dirs[name]; exists {
		log.Infof("%s entries %d", name, len(d.Entries))
		return &d, nil
	}
	return nil, os.ErrNotExist
}

var _ os.FileInfo = (*FileInfo)(nil)

// the awkward File* prefix is to export the field but avoid conflict with os.FileInfo interface ...
type FileInfo struct {
	FileName    string
	FileSize    int64
	FileMode    os.FileMode
	FileModTime time.Time
	FileIsDir   bool
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

func init() {
	registeredBoxes = make(map[string]EmbedBox)
}
