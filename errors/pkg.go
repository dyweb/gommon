// Package errors provides multi error, error wrapping. It defines error category code for machine post processing
package errors

import "fmt"

const (
	// MultiErrSep is the separator used when returning a slice of errors as single line message
	MultiErrSep = "; "
	// ErrCauseSep is the separator used when returning a error with causal chain as single line message
	ErrCauseSep = ": "
)

func New(msg string) error {
	return &FreshError{
		msg:   msg,
		stack: callers(),
	}
}

func Errorf(format string, args ...interface{}) error {
	return &FreshError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}
