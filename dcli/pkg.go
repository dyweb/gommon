// Package dcli is a commandline application builder.
// It supports git style sub command and is modeled after spf13/cobra.
package dcli

import (
	dlog "github.com/dyweb/gommon/log"
)

var logReg = dlog.NewRegistry()
var log = logReg.NewLogger()
