// Package linter use static analysis to enforce customized coding style and fix(format) go code.
package linter

import (
	dlog "github.com/dyweb/gommon/log"
)

var (
	logReg = dlog.NewRegistry()
	log    = logReg.Logger()
)

// Level describes lint error level
// TODO: sync w/ existing lint tools
type Level int

const (
	UnknownLevel Level = iota
	Deprecated
	Warn
	Error
)
