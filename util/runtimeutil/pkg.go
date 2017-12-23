package runtimeutil

import (
	"runtime"
	"strings"
)

// GetCallerPackage is used by log package to get caller source code position
func GetCallerPackage(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	// FIXME: https://github.com/golang/go/issues/19426 use runtime.Frames instead of runtime.FuncForPC
	fn := runtime.FuncForPC(pc)
	fnName := fn.Name()
	lastDot := strings.LastIndex(fnName, ".")
	return fnName[:lastDot]
}

// TODO: we can get file and line as well ...
//func GetCallerPackageAndFunc(skip int) (string, string) {
//
//}
