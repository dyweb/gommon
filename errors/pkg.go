// Package errors provides multi error, wrapping and inspection.
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

// Ignore swallow the error, you should NOT use it unless you know what you are doing (make the lint tool happy)
// It is inspired by dgraph x/error.go
func Ignore(_ error) {
	// do nothing
}

// Ignore2 ignores return value and error, it is useful for functions like Write(b []byte) (int64, error)
// It is also inspired by dgraph x/error.go
func Ignore2(_ interface{}, _ error) {
	// do nothing
}

// yeah, there is no Ignore3
