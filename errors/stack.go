package errors

import (
	"fmt"
	"runtime"
)

// TODO: handling print stack

func callers() *runtime.Frames {
	// TODO: it should be configurable? or move to pkg.go?
	const depth = 20
	var pcs = make([]uintptr, depth)
	// 3 skips runtime.Callers itself, callers function, and the function that creates error, i.e New, Errorf
	//more true | runtime.Callers /home/at15/app/go/src/runtime/extern.go:212
	//more true | github.com/dyweb/gommon/errors.callers /home/at15/workspace/src/github.com/dyweb/gommon/errors/stack.go:18
	//more true | github.com/dyweb/gommon/errors.New /home/at15/workspace/src/github.com/dyweb/gommon/errors/pkg.go:16
	n := runtime.Callers(3, pcs)
	return runtime.CallersFrames(pcs[:n])
}

func printFrames(frames *runtime.Frames) {
	// from https://golang.org/pkg/runtime/#Frames
	for {
		frame, more := frames.Next()
		fmt.Printf("more %v | %s %s:%d\n", more, frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func framesLen(frames *runtime.Frames) int {
	i := 0
	for {
		_, more := frames.Next()
		if more {
			i++
		} else {
			return i
		}
	}
}
