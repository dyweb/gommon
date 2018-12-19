// Package errors provides error wrapping, multi error and error inspection.
package errors

import (
	"fmt"
)

const (
	// MultiErrSep is the separator used when returning a slice of errors as single line message
	MultiErrSep = "; "
	// ErrCauseSep is the separator used when returning a error with causal chain as single line message
	ErrCauseSep = ": "
)

// New creates a freshError with stack trace
func New(msg string) error {
	return &freshError{
		msg:   msg,
		stack: callers(),
	}
}

// Errorf is New with fmt.Sprintf
func Errorf(format string, args ...interface{}) error {
	return &freshError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}
