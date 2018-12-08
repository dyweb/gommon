package log

import "runtime"

// caller.go contains logic for getting log file line

// Caller is the result from runtime.Caller
type Caller struct {
	// File is the full file path without any trimming, it would be UnknownFile if caller is not found
	File string
	Line int
}

var emptyCaller = Caller{File: "", Line: 0}

// EmptyCaller is mainly used for testing handler, it contains a empty file and line 0
func EmptyCaller() Caller {
	return emptyCaller
}

const UnknownFile = "<?>"

// caller gets source location at runtime, in the future we may generate it at compile time to reduce the
// overhead, though I am not sure what the overhead is without actual benchmark and profiling
// Also I think the complexity added does not worth the performance benefits it gives
// TODO: https://github.com/dyweb/gommon/issues/43
// TODO: add test for skip ...
func caller(skip int) Caller {
	_, file, line, ok := runtime.Caller(2 + skip)
	if !ok {
		return Caller{UnknownFile, 1}
	}
	return Caller{File: file, Line: line}
}
