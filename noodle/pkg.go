// Package noodle helps embedding static assets into go binary, it supports using ignore file
package noodle // import "github.com/dyweb/gommon/noodle"

import (
	"github.com/dyweb/gommon/util/logutil"
	"net/http"
)

const (
	DefaultIgnoreFileName = ".noodleignore"
)

var log = logutil.NewPackageLogger()

// Bowel is the container for different types of noodles
type Bowel interface {
	http.FileSystem
}
