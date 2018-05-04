// Package logutil is a registry of loggers, it is required for all lib and app that use gommon/log.
// You should add the registry as child of your library/application's child if you want to control gommon libraries
// logging behavior
package logutil // import "github.com/dyweb/gommon/util/logutil"

import (
	"github.com/dyweb/gommon/log"
)

var Registry = log.NewLibraryLogger()

func NewPackageLogger() *log.Logger {
	l := log.NewPackageLoggerWithSkip(1)
	Registry.AddChild(l)
	return l
}
