package errors

import (
	"fmt"
	"runtime"
)

// TODO: handling print stack

func callers() []runtime.Frame {
	// TODO: it should be configurable? or move to pkg.go?
	const depth = 20
	pcs := make([]uintptr, depth)
	// 3 skips runtime.Callers itself, callers function, and the function that creates error, i.e New, Errorf
	//more true | runtime.Callers /home/at15/app/go/src/runtime/extern.go:212
	//more true | github.com/dyweb/gommon/errors.callers /home/at15/workspace/src/github.com/dyweb/gommon/errors/stack.go:18
	//more true | github.com/dyweb/gommon/errors.New /home/at15/workspace/src/github.com/dyweb/gommon/errors/pkg.go:16
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	res := make([]runtime.Frame, 0, depth)
	for {
		frame, more := frames.Next()
		res = append(res, frame)
		if !more {
			break
		}
	}
	return res
}

func printFrames(frames []runtime.Frame) {
	for i := 0; i < len(frames); i++ {
		fmt.Printf("%s %s:%d\n", frames[i].Function, frames[i].File, frames[i].Line)
	}
}

func printFramesPtr(frames *runtime.Frames) {
	// from https://golang.org/pkg/runtime/#Frames
	for {
		frame, more := frames.Next()
		fmt.Printf("more %v | %s %s:%d\n", more, frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}
