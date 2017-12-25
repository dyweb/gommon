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

// see https://github.com/dyweb/gommon/issues/32
// based on https://github.com/go-stack/stack/blob/master/stack.go#L29:51
// TODO: not sure if calling two Next without checking the more value works for other go version
func GetCallerFrame(skip int) runtime.Frame {
	var pcs [3]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	f, _ := frames.Next()
	f, _ = frames.Next()
	return f
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

func SplitStructMethod(f string) (st string, function string) {
	dot := strings.LastIndex(f, ".")
	st, function = f[:dot], f[dot+1:]
	if st[0] == '(' {
		st = st[1:len(st)-1]
	}
	if st[0] == '*' {
		st = st[1:]
	}
	return
}
