package noodle

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

var registeredBoxes map[string]EmbedBox

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

var _ os.FileInfo = (*FileInfo)(nil)

// the awkward File* prefix is to export the field but avoid conflict with os.FileInfo interface ...
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

func init() {
	registeredBoxes = make(map[string]EmbedBox)
}
