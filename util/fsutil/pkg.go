// Package fsutil adds ignore support for walk
package fsutil // import "github.com/dyweb/gommon/util/fsutil"

import (
	dlog "github.com/dyweb/gommon/log"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()
