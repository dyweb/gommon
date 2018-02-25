package errors

import (
	"fmt"
	"runtime"
)

// TODO: it should be configurable? or move to pkg.go?
const depth = 20

type Stack struct {
	p      *runtime.Frames
	depth  int
	frames []runtime.Frame
}

func (s *Stack) Frames() []runtime.Frame {
	if s == nil || s.p == nil {
		return nil
	}
	if len(s.frames) != 0 {
		return s.frames
	}
	frames := make([]runtime.Frame, 0, s.depth)
	for {
		frame, more := s.p.Next()
		frames = append(frames, frame)
		if !more {
			break
		}
	}
	s.frames = frames
	return frames
}

// TODO: handling print stack

func callers() *Stack {
	pcs := make([]uintptr, depth)
	// 3 skips runtime.Callers itself, callers function, and the function that creates error, i.e New, Errorf
	//more true | runtime.Callers /home/at15/app/go/src/runtime/extern.go:212
	//more true | github.com/dyweb/gommon/errors.callers /home/at15/workspace/src/github.com/dyweb/gommon/errors/stack.go:18
	//more true | github.com/dyweb/gommon/errors.New /home/at15/workspace/src/github.com/dyweb/gommon/errors/pkg.go:16
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	return &Stack{
		p:     frames,
		depth: n,
	}
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
