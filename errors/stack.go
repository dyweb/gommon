package errors

import (
	"fmt"
	"runtime"
)

// TODO: it should be configurable? or move to pkg.go and make it public?
// TODO: what if there is recursive function .... I guess any depth is not enough then? what is the limit on go runtime's side?
const depth = 20

type Stack interface {
	Frames() []runtime.Frame
}

// lazyStack does not loop the frames until its Frames method is called,
// It will the existing slice of frames if presented, so you can use a Stack interface to create a lazyStack
// s := lazyStack{frames: e.(Stack).Frames()}
type lazyStack struct {
	p      *runtime.Frames
	depth  int
	frames []runtime.Frame
}

func (s *lazyStack) Frames() []runtime.Frame {
	if s == nil {
		return nil
	}
	// we can init lazyStack using exist frames
	if len(s.frames) != 0 {
		return s.frames
	}
	// no existing stack and no pointer to frames, can't move on
	if s.p == nil {
		return nil
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

func callers() *lazyStack {
	pcs := make([]uintptr, depth)
	// 3 skips runtime.Callers itself, callers function, and the function that creates error, i.e New, Errorf
	//more true | runtime.Callers /home/at15/app/go/src/runtime/extern.go:212
	//more true | github.com/dyweb/gommon/errors.callers /home/at15/workspace/src/github.com/dyweb/gommon/errors/stack.go:18
	//more true | github.com/dyweb/gommon/errors.New /home/at15/workspace/src/github.com/dyweb/gommon/errors/pkg.go:16
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	return &lazyStack{
		p:     frames,
		depth: n,
	}
}

func PrintFrames(frames []runtime.Frame) {
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
