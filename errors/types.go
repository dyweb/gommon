package errors

import "runtime"

// TODO: error types, user, dev, std, lib, grpc etc.

func IsRuntimeError(err error) bool {
	_, ok := err.(runtime.Error)
	return ok
}
