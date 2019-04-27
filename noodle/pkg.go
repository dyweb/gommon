// Package noodle helps embedding static assets into go binary, it supports using ignore file
package noodle

import (
	"net/http"

	dlog "github.com/dyweb/gommon/log"
)

const (
	DefaultIgnoreFileName = ".noodleignore"
	DefaultName           = "Bowel"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()

// Bowel is the container for different types of noodles
type Bowel interface {
	http.FileSystem
}
