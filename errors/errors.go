package errors

import "fmt"

// Causer return the underlying error, a error does not have cause should return itself.
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

type TracedError interface {
	fmt.Formatter
	ErrorStack() *Stack
}

var (
	_ error       = (*FreshError)(nil)
	_ TracedError = (*FreshError)(nil)
	_ Causer      = (*FreshError)(nil)
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

func (fresh *FreshError) Cause() error {
	return fresh
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

// Wrap attach stack to
func Wrap(err error, msg string) error {
	// NOTE: sometimes we call wrap without check if the error is nil, it is cleaner if it is the last statement in func
	//
	// i.e. return errors.Wrap(f.Close(), "failed to close file")
	//
	// 		if err := f.Close(); err != nil {
	//			return errors.Wrap(err, "failed to close file")
	//      }
	//      return nil
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

func Wrapf(err error, format string, args ...interface{}) error {
	// NOTE: copied from wrap instead of call Wrap due to caller
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
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: stack,
	}
}

// Cause returns cause of the error, it stops at the last error that does not implement Causer interface
// errors wrapped using pkg/errors also satisfy and this interface can be unwrapped as well.
// If error is nil, it will return nil.
// If error is a standard library error or FreshError it will return the error itself.
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

func (wrapped *WrappedError) Error() string {
	return wrapped.msg + ": " + wrapped.cause.Error()
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
