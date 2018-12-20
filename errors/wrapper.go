package errors

import (
	"fmt"
	"reflect"
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
	// NOTE: keep the following in sync with Wrapf
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

// Is walks the error chain and do direct compare.
// It should be used for checking sentinel errors like io.EOF
// It returns true on the first match.
// It returns false when there is no match.
//
// It unwraps both wrapped error and multi error
func Is(err, target error) bool {
	for {
		// TODO: put this before nil check allow Is(nil, nil) returns true, not sure if this is desired behaviour
		if err == target {
			return true
		}
		if err == nil {
			return false
		}
		// support both Causer from pkg/errors and Wrapper from go 2 proposal
		switch err.(type) {
		case Wrapper:
			err = err.(Wrapper).Unwrap()
		case causer:
			err = err.(causer).Cause()
		case ErrorList:
			// support multi error
			errs := err.(ErrorList).Errors()
			for i := 0; i < len(errs); i++ {
				if Is(errs[i], target) {
					return true
				}
			}
			return false // NOTE: without the return we will be looping forever because err is not updated
		default:
			return false
		}
	}
}

// IsType walks the error chain and match by type using reflect.
// It only returns match result, if you need the matched error, use GetType
//
// It should be used for checking errors that defined their own types,
// errors created using errors.New, errors.Errof should NOT be checked using this method
// because they have same type, string if you are using standard library, freshError if you are using gommon/errors
//
// It calls IsTypeOf to reduce the overhead of calling reflect on target error
func IsType(err, target error) bool {
	_, ok := GetType(err, target)
	return ok
}

// IsTypeOf requires user to call reflect.TypeOf(exampleErr).String() as the type string
func IsTypeOf(err error, tpe string) bool {
	_, ok := GetTypeOf(err, tpe)
	return ok
}

// GetType walks the error chain and match by type using reflect,
// It returns the matched error and match result.
// You still need to do a type conversion on the returned error.
//
// It calls GetTypeOf to reduce the overhead of calling reflect on target error
func GetType(err, target error) (matched error, ok bool) {
	if err == nil || target == nil {
		return nil, false
	}
	return GetTypeOf(err, reflect.TypeOf(target).String())
}

// GetTypeOf requires user to call reflect.TypeOf(exampleErr).String() as the type string
func GetTypeOf(err error, tpe string) (error, bool) {
	if err == nil {
		return nil, false
	}
	for {
		if err == nil {
			return nil, false
		}
		if reflect.TypeOf(err).String() == tpe {
			return err, true
		}
		switch err.(type) {
		case Wrapper:
			err = err.(Wrapper).Unwrap()
		case causer:
			err = err.(causer).Cause()
		case ErrorList:
			errs := err.(ErrorList).Errors()
			for i := 0; i < len(errs); i++ {
				m, ok := GetTypeOf(errs[i], tpe)
				if ok {
					return m, true
				}
			}
			return nil, false
		default:
			return nil, false
		}
	}
}

// AsValue is in go 2 proposal as workaround if go 2 does not have polymorphism,
// however, it's pretty hard to use and user can have error easily,
// we decided to use GetType
//func AsValue(val interface{}, err error) bool {
//
//}

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
