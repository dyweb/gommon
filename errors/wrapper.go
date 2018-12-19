package errors

import (
	"fmt"
)

// wrapper.go defines interface and util method for error wrapping

// Wrapper is based on go 2 proposal, it only has an Unwrap method to returns the underlying error
type Wrapper interface {
	Unwrap() error
}

// Message return the top level error message without traverse the error chain
// i.e. when Error() returns `invalid config: file a.json does not exist` Message() returns `invalid config`
type Messenger interface {
	Message() string
}

// Tracer is error with stack trace
type Tracer interface {
	Stack() Stack
}

// Wrap creates a wrappedError with stack and set its cause to err.
//
// If the error being wrapped is already a Tracer, Wrap will reuse its stack trace instead of creating a new one.
// The error being wrapped has deeper stack than where the Wrap function is called and is closer to the root of error.
// This is based on https://github.com/pkg/errors/pull/122 to avoid having extra interface like WithMessage and WithStack
// like https://github.com/pkg/errors does.
//
// Wrap returns nil if the error you are trying to wrap is nil, thus if it is the last error checking in a func,
// you can return the wrap function directly in one line instead of using typical three line error check and wrap.
//
//      return errors.Wrap(f.Close(), "failed to close file")
//
//      if err := f.Close(); err != nil {
//            return errors.Wrap(err, "failed to close file")
//      }
//      return nil。。
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if err == nil {
		return nil
	}
	var stack *lazyStack
	// reuse existing stack
	if t, ok := err.(Tracer); ok {
		stack = &lazyStack{
			frames: t.Stack().Frames(),
		}
	} else {
		stack = callers()
	}
	return &wrappedError{
		msg:   msg,
		cause: err,
		stack: stack,
	}
}

// Wrapf is Wrap with fmt.Sprintf
func Wrapf(err error, format string, args ...interface{}) error {
	// NOTE: copied from wrap instead of call Wrap due to caller
	// -- copy & paste start
	if err == nil {
		return nil
	}
	var stack *lazyStack
	// reuse existing stack
	if t, ok := err.(Tracer); ok {
		stack = &lazyStack{
			frames: t.Stack().Frames(),
		}
	} else {
		stack = callers()
	}
	// --- copy & paste end
	return &wrappedError{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: stack,
	}
}

var (
	_ error         = (*freshError)(nil)
	_ Tracer        = (*freshError)(nil)
	_ fmt.Formatter = (*freshError)(nil)
)

// freshError is a root error with stack trace
type freshError struct {
	msg   string
	stack *lazyStack
}

func (fresh *freshError) Error() string {
	return fresh.msg
}

func (fresh *freshError) Stack() Stack {
	return fresh.stack
}

func (fresh *freshError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// TODO: print stack if s.Flag('+')
		fallthrough
	case 's':
		Ignore2(s.Write([]byte(fresh.msg)))
	case 'q':
		// %q	a double-quoted string safely escaped with Go syntax
		Ignore2(fmt.Fprintf(s, "%q", fresh.msg))
	}
}

var (
	_ error         = (*wrappedError)(nil)
	_ Tracer        = (*wrappedError)(nil)
	_ causer        = (*wrappedError)(nil)
	_ Wrapper       = (*wrappedError)(nil)
	_ fmt.Formatter = (*wrappedError)(nil)
)

// wrappedError implements the Wrapper
type wrappedError struct {
	msg   string
	cause error
	stack *lazyStack
}

func (wrapped *wrappedError) Error() string {
	return wrapped.msg + ErrCauseSep + wrapped.cause.Error()
}

func (wrapped *wrappedError) Stack() Stack {
	return wrapped.stack
}

// Deprecated: use Unwrap
func (wrapped *wrappedError) Cause() error {
	return wrapped.cause
}

func (wrapped *wrappedError) Unwrap() error {
	return wrapped.cause
}

func (wrapped *wrappedError) Message() string {
	return wrapped.msg
}

func (wrapped *wrappedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// TODO: print stack if s.Flag('+')
		fallthrough
	case 's':
		Ignore2(s.Write([]byte(wrapped.Error())))
	case 'q':
		Ignore2(fmt.Fprintf(s, "%q", wrapped.Error()))
	}
}
