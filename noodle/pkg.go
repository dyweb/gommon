package noodle

import (
	"github.com/dyweb/gommon/util/logutil"
	"path/filepath"
)

const (
	ignoreFileName = ".noodleignore"
)

var log = logutil.NewPackageLogger()

func join(elem ...string) string {
	return filepath.Join(elem...)
}
