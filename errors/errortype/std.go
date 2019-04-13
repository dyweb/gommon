package errortype

import (
	"net"
	"os"
	"runtime"

	"github.com/dyweb/gommon/errors"
)

// NOTE: this is not the file that defines a sort of structs ....
// TODO: error types, user, dev, std, lib, grpc etc. see https://github.com/dyweb/gommon/issues/66

func IsRuntimeError(err error) bool {
	_, ok := GetRuntimeError(err)
	return ok
}

// RuntimeError is triggered by invalid operation like divide by 0, index out of range in slice.
// When user call panic, which accepts an interface, most of them will just give a simple string
// and does not satisfies the runtime.Error interface
func GetRuntimeError(err error) (rErr runtime.Error, ok bool) {
	errors.Walk(err, func(err error) (stop bool) {
		rErr, ok = err.(runtime.Error)
		if ok {
			return true
		}
		return false
	})
	return
}

func IsNetError(err error) bool {
	_, ok := GetNetError(err)
	return ok
}

func GetNetError(err error) (netErr net.Error, ok bool) {
	errors.Walk(err, func(err error) (stop bool) {
		netErr, ok = err.(net.Error)
		if ok {
			return true // stop
		}
		return false // continue
	})
	return
}

func IsFsError(err error) bool {
	_, ok := GetFsError(err)
	return ok
}

func GetFsError(err error) (error, bool) {
	var found error
	errors.Walk(err, func(err error) (stop bool) {
		// copied from os/error.go underlyingError
		switch err := err.(type) {
		case *os.PathError:
			found = err
			return true
		case *os.LinkError:
			found = err
			return true
		case *os.SyscallError:
			found = err
			return true
		}
		return
	})
	if found != nil {
		return found, true
	}
	return found, false
}
