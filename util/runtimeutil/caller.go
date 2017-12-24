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

// SplitPackageFunc returns package (w/o GOPATH) and function (w/ struct if presented) based on runtime.Frame.Function
// Copied from runtime.Frame struct documentation
// Func may be nil for non-Go code or fully inlined functions
// If Func is not nil then Function == Func.Name()
// github.com/dyweb/gommon/log2/_examples/uselib/service.(*Auth).Check
// github.com/dyweb/gommon/log2.TestNewIdentityFromCaller
func SplitPackageFunc(f string) (pkg string, function string) {
	dot := 0
	// go from back of the string
	// the first dot splits package (w/ struct) and function, the second dot split package and struct (if any)
	// we put struct (if any) and function together, so we just need to dot closest to last /
	for i := len(f) - 1; i >= 0; i-- {
		// TODO: it might not work on windows
		if f[i] == '/' {
			break
		}
		if f[i] == '.' {
			dot = i
		}
	}
	return f[:dot], f[dot+1:]
}
