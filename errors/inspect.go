package errors

import "reflect"

// inspect.go defines functions for inspecting wrapped error or error list

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

// WalkFunc accepts an error and based on its inspection logic it can tell
// Walk if it should stop walking the error chain or error list
type WalkFunc func(err error) (stop bool)

// Walk traverse error chain and error list, it stops when there is no
// underlying error or the WalkFunc decides to stop
// TODO: might let Is and GetType use Walk, this reduce copy and paste ...
func Walk(err error, cb WalkFunc) {
	if err == nil {
		return
	}
	for {
		if err == nil {
			return
		}
		// WalkFunc decides to stop
		if cb(err) {
			return
		}
		switch err.(type) {
		case Wrapper:
			err = err.(Wrapper).Unwrap()
		case causer:
			err = err.(causer).Cause()
		case ErrorList:
			errs := err.(ErrorList).Errors()
			for i := 0; i < len(errs); i++ {
				Walk(errs[i], cb)
			}
			return
		default:
			return
		}
	}
}
