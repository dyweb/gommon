// Package noodle helps embedding static assets into go binary, it supports using ignore file
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
