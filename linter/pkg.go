// Package linter use static analysis to enforce customized coding style and fix(format) go code.
package linter

import (
	dlog "github.com/dyweb/gommon/log"
)

var (
	logReg = dlog.NewRegistry()
	log    = logReg.Logger()
)
