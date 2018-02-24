// Package logutil is a registry of loggers, it is required for all lib and app that use gommon/log
package logutil

import (
	"github.com/dyweb/gommon/log"
)

var Registry = log.NewLibraryLogger()

func NewPackageLogger() *log.Logger {
	l := log.NewPackageLoggerWithSkip(1)
	Registry.AddChild(l)
	return l
}
