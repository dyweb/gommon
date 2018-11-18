// Package noodle helps embedding static assets into go binary, it supports using ignore file
package noodle

import (
	"net/http"

	"github.com/dyweb/gommon/util/logutil"
)

const (
	DefaultIgnoreFileName = ".noodleignore"
	DefaultName           = "Bowel"
)

var log = logutil.NewPackageLogger()

// Bowel is the container for different types of noodles
type Bowel interface {
	http.FileSystem
}
