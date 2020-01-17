# juju/errors

- repo: https://github.com/juju/errors
- doc: https://godoc.org/github.com/juju/errors
- last update: 2017/06/02

- it is not using errgo ... https://github.com/juju/errors/pull/10
  - the description of repo is misleading ...
  
  
Usage

- mostly same as juju-errgo
- location is ~~not set when error is created~~ when using `Wrapf` etc. location is set
- [ ] its way of obtain go path may not work when `vendor` is used ...

````go
// Err holds a description of an error along with information about
// where the error was created.
//
// It may be embedded in custom error types to add extra information that
// this errors package can understand.
type Err struct {
	// message holds an annotation of the error.
	message string

	// cause holds the cause of the error as returned
	// by the Cause method.
	cause error

	// previous holds the previous error in the error stack, if any.
	previous error

	// file and line hold the source code location where the error was
	// created.
	file string
	line int
}


// Cause returns the cause of the given error.  This will be either the
// original error, or the result of a Wrap or Mask call.
//
// Cause is the usual way to diagnose errors that may have been wrapped by
// the other errors functions.
func Cause(err error) error {
	var diag error
	if err, ok := err.(causer); ok {
		diag = err.Cause()
	}
	if diag != nil {
		return diag
	}
	return err
}

type causer interface {
	Cause() error
}

type wrapper interface {
	// Message returns the top level error message,
	// not including the message from the Previous
	// error.
	Message() string

	// Underlying returns the Previous error, or nil
	// if there is none.
	Underlying() error
}
````