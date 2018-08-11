package errors

import "fmt"

// Causer returns the underlying error, a error without cause should return itself.
// It is based on the private `causer` interface in pkg/errors, so errors wrapped using pkg/errors can also be handled
type Causer interface {
	Cause() error
}

// Wrapper has cause and its own error message. It is based on the private `wrapper` interface in juju/errors
type Wrapper interface {
	Causer
	// Message return the top level error message without concat message from its cause
	// i.e. when Error() returns `invalid config: file a.json does not exist` Message() returns `invalid config`
	Message() string
}

// TracedError is error with stack trace
type TracedError interface {
	fmt.Formatter
	ErrorStack() *Stack
}

// Wrap creates a WrappedError with stack and set its cause to err.
//
// If the error being wrapped is already a TracedError, Wrap will reuse its stack trace instead of creating a new one.
// The error being wrapped has deeper stack than where the Wrap function is called and is closer to the root of error.
// This is based on https://github.com/pkg/errors/pull/122 to avoid having extra interface like WithMessage and WithStack
// like https://github.com/pkg/errors does.
//
// Wrap returns nil if the error you are trying to wrap is nil, thus if it is the last error checking, you can return
// the wrap function directly in one line instead of using typical three line error check and wrap. i.e.
//
//      return errors.Wrap(f.Close(), "failed to close file")
//
//      if err := f.Close(); err != nil {
//            return errors.Wrap(err, "failed to close file")
//      }
//      return nil
//
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	var stack *Stack
	// reuse existing stack
	if t, ok := err.(TracedError); ok {
		stack = t.ErrorStack()
	} else {
		stack = callers()
	}
	return &WrappedError{
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
	var stack *Stack
	// reuse existing stack
	if t, ok := err.(TracedError); ok {
		stack = t.ErrorStack()
	} else {
		stack = callers()
	}
	// --- copy & paste end
	return &WrappedError{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: stack,
	}
}

// Cause returns root cause of the error (if any), it stops at the last error that does not implement Causer interface.
// If you want get direct cause, use DirectCause.
// If error is nil, it will return nil. If error is not wrapped it will return the error itself.
// error wrapped using https://github.com/pkg/errors also satisfies this interface and can be unwrapped as well.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	for err != nil {
		causer, ok := err.(Causer)
		if !ok {
			break
		}
		err = causer.Cause()
	}
	return err
}

// DirectCause returns the direct cause of the error (if any). It does NOT follow the cause chain, if you want to get
// root cause, use Cause
func DirectCause(err error) error {
	if err == nil {
		return nil
	}
	causer, ok := err.(Causer)
	if !ok {
		return nil
	}
	return causer.Cause()
}

var (
	_ error       = (*FreshError)(nil)
	_ TracedError = (*FreshError)(nil)
)

type FreshError struct {
	msg   string
	stack *Stack
}

func (fresh *FreshError) Error() string {
	return fresh.msg
}

func (fresh *FreshError) ErrorStack() *Stack {
	return fresh.stack
}

func (fresh *FreshError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// TODO: print stack if s.Flag('+')
		fallthrough
	case 's':
		s.Write([]byte(fresh.msg))
	case 'q':
		// %q	a double-quoted string safely escaped with Go syntax
		fmt.Fprintf(s, "%q", fresh.msg)
	}
}

var (
	_ error       = (*WrappedError)(nil)
	_ TracedError = (*WrappedError)(nil)
	_ Causer      = (*WrappedError)(nil)
	_ Wrapper     = (*WrappedError)(nil)
)

type WrappedError struct {
	msg   string
	cause error
	stack *Stack
}

func (wrapped *WrappedError) Error() string {
	return wrapped.msg + ErrCauseSep + wrapped.cause.Error()
}

func (wrapped *WrappedError) ErrorStack() *Stack {
	return wrapped.stack
}

func (wrapped *WrappedError) Cause() error {
	return wrapped.cause
}

func (wrapped *WrappedError) Message() string {
	return wrapped.msg
}

func (wrapped *WrappedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// TODO: print stack if s.Flag('+')
		fallthrough
	case 's':
		s.Write([]byte(wrapped.Error()))
	case 'q':
		fmt.Fprintf(s, "%q", wrapped.Error())
	}
}
