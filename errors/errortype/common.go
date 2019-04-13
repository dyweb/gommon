package errortype

import "github.com/dyweb/gommon/errors"

// common.go defines common error types factory functions

type ErrNotFound interface {
	IsErrNotFound()
}

func IsNotFound(err error) bool {
	is := false
	// TODO: can I use and interface as argument
	errors.Walk(err, func(err error) (stop bool) {
		_, ok := err.(ErrNotFound)
		if ok {
			is = true
			return false
		}
		return false
	})
	return is
}

type errNotFound struct {
	lookingFor string
}

func (e *errNotFound) IsErrNotFound() {}

func (e *errNotFound) Error() string {
	return e.lookingFor + " not found"
}

func NewNotFound(lookingFor string) error {
	return &errNotFound{lookingFor: lookingFor}
}

type ErrAlreadyExists interface {
	IsErrAlreadyExists()
}

// FIXME: this is copy and paste, but can't pass an interface directly
func IsAlreadyExists(err error) bool {
	is := false
	// TODO: can I use and interface as argument
	errors.Walk(err, func(err error) (stop bool) {
		_, ok := err.(ErrAlreadyExists)
		if ok {
			is = true
			return false
		}
		return false
	})
	return is
}

type errAlreadyExists struct {
	name string
}

func (e *errAlreadyExists) IsErrAlreadyExists() {}

func (e *errAlreadyExists) Error() string {
	return e.name + " already exists"
}

func NewAlreadyExists(name string) error {
	return &errAlreadyExists{name: name}
}

type ErrNotImplemented interface {
	IsErrNotImplemented()
}

// FIXME: this is copy and paste, but can't pass an interface directly
func IsNotImplemented(err error) bool {
	is := false
	// TODO: can I use and interface as argument
	errors.Walk(err, func(err error) (stop bool) {
		_, ok := err.(ErrNotImplemented)
		if ok {
			is = true
			return false
		}
		return false
	})
	return is
}

type errNotImplemented struct {
	feature string
}

func (e *errNotImplemented) IsErrNotImplemented() {}

func (e *errNotImplemented) Error() string {
	return e.feature + " not implemented"
}

func NewNotImplemented(feature string) error {
	return &errNotImplemented{feature: feature}
}
