package errors

import "fmt"

type TracedError interface {
	fmt.Formatter
	ErrorStack() *Stack
}

var _ error = (*FreshError)(nil)
var _ TracedError = (*FreshError)(nil)

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

var _ error = (*WrappedError)(nil)
var _ TracedError = (*WrappedError)(nil)

type WrappedError struct {
	msg   string
	cause error
	stack *Stack
}

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

func (wrapped *WrappedError) Error() string {
	return wrapped.msg + ": " + wrapped.cause.Error()
}

func (wrapped *WrappedError) ErrorStack() *Stack {
	return wrapped.stack
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
