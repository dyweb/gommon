package errors

import (
	"reflect"
)

// inspect.go defines functions for inspecting wrapped error or error list

// Is walks the error chain and do direct compare.
// It should be used for checking sentinel errors like io.EOF
// It returns true on the first match.
// It returns false when there is no match.
//
// It unwraps both wrapped error and multi error.
func Is(err, target error) bool {
	// NOTE: in std it returns err == target when target == nil
	if err == nil || target == nil {
		return err == target
	}

	found := false
	Walk(err, func(err error) (stop bool) {
		// TODO: in std it checks if err and target are comparable
		// TODO: std supports Is method on err to define custom equality logic
		if err == target {
			found = true
			return true
		}
		return false
	})
	return found
}

// IsType walks the error chain and match by type using reflect.
// It only returns match result, if you need the matched error, use GetType
//
// It should be used for checking errors that defined their own types.
// errors created using errors.New, errors.Errof should NOT be checked using this method
// because they have same type, string if you are using standard library, freshError if you are using gommon/errors
//
// It calls IsTypeOf to reduce the overhead of calling reflect on target error
func IsType(err, target error) bool {
	_, ok := GetType(err, target)
	return ok
}

// IsTypeOf requires user to call reflect.TypeOf(exampleErr).String() as the type string
func IsTypeOf(err error, tpe reflect.Type) bool {
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
	return GetTypeOf(err, reflect.TypeOf(target))
}

// GetTypeOf requires user to call reflect.TypeOf(exampleErr) as target type.
// NOTE: for the type matching, we compare equality of reflect.Type directly,
// Originally we were comparing string, however result from `String()` is not
// the full path, so x/foo/bar.Encoder can match y/foo/bar.Encoder because we
// got bar.Encoder for both of them.
// You can compare interface{} if their underlying type is same and comparable,
// it is documented in https://golang.org/pkg/reflect/#Type
//
// Related https://github.com/dyweb/gommon/issues/104
func GetTypeOf(err error, tpe reflect.Type) (matched error, ok bool) {
	if err == nil {
		return nil, false
	}

	Walk(err, func(err error) (stop bool) {
		if reflect.TypeOf(err) == tpe {
			matched = err
			return true
		}
		return false
	})
	return matched, matched != nil
}

// As walks the error chain and save the matched error in target.
// target should be a pointer to an error type.
func As(err error, target interface{}) bool {
	if err == nil || target == nil {
		// TODO: std panic when target is nil
		return false
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}
	targetType := val.Type().Elem()

	found := false
	Walk(err, func(err error) (stop bool) {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			found = true
			return true
		}
		return false
	})
	return found
}

// WalkFunc accepts an error and based on its inspection logic it can tell
// Walk if it should stop walking the error chain or error list
type WalkFunc func(err error) (stop bool)

// Walk traverse error chain and error list, it stops when there is no
// underlying error or the WalkFunc decides to stop.
func Walk(err error, cb WalkFunc) {
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
