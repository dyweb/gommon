// Package fsutil adds ignore support for walk
package fsutil

import (
	dlog "github.com/dyweb/gommon/log"
)

var (
	logReg = dlog.NewRegistry()
	log    = logReg.Logger()
)
