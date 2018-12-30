package errortype

import "runtime"

// NOTE: this is not the file that defines a sort of structs ....
// TODO: error types, user, dev, std, lib, grpc etc. see https://github.com/dyweb/gommon/issues/66
// TODO: Is function should do the unwrapping as well
// TODO: should consider put it into another package like errort, now it's errortype

func IsRuntimeError(err error) bool {
	_, ok := err.(runtime.Error)
	return ok
}

