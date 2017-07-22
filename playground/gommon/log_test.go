package gommon

import (
	"testing"
	dlog "github.com/dyweb/gommon/log"
)

var Logger = dlog.NewLogger()
var log = Logger.RegisterPkg()

func TestLogger_RegisterPkg(t *testing.T) {
	log.Info("should show the right package")
}

func init() {
	f := dlog.NewTextFormatter()
	f.EnableColor = true
	Logger.Formatter = f
	Logger.Level = dlog.InfoLevel
}
