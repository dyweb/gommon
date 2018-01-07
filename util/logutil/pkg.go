package logutil

import (
	"github.com/dyweb/gommon/log"
)

var Logger = log.NewLibraryLogger()

func NewLogger() *log.Logger {
	l := log.NewPackageLoggerWithSkip(1)
	Logger.AddChild(l)
	return l
}

func init() {

}
