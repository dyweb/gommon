package runtimeutil

import (
	"runtime"
	"strings"
	"fmt"
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

// TODO: doc
// Copied from runtime.Frame struct documentation
// Func may be nil for non-Go code or fully inlined functions
// If Func is not nil then Function == Func.Name()
// github.com/dyweb/gommon/log2/_examples/uselib/service.(*Auth).Check
// github.com/dyweb/gommon/log2.TestNewIdentityFromCaller
func SplitPackageFunc(f string) (pkg string, function string) {
	// go from back of the string, the first dot splits package (w/ struct) and function
	funcDot := strings.LastIndex(f, ".")
	fmt.Println("funcDot", funcDot)
	pkgWithStruct := f[:funcDot]
	fmt.Println("pkgWithStruct", pkgWithStruct)
	// the second dot splits package and struct (if there is any)
	structDot := strings.LastIndex(pkgWithStruct, ".")
	// FIXME: this won't work, it will find the dot in github.com ....
	if structDot == -1 {
		return pkgWithStruct, f[funcDot+1:]
	}
	return f[:structDot], f[structDot+1:]
}
