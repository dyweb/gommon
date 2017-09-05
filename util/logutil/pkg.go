package logutil

import (
	"github.com/dyweb/gommon/log"
)

// Logger is the default logger with info level
var Logger = log.NewLogger()

func init() {
	f := log.NewTextFormatter()
	f.EnableColor = true
	Logger.Formatter = f
	Logger.Level = log.InfoLevel
}
