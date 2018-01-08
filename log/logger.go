package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	mu       sync.RWMutex
	h        Handler
	level    Level
	fields   Fields
	children map[uint64][]*Logger
	id       *Identity
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

func (l *Logger) SetHandler(h Handler) {
	l.mu.Lock()
	l.h = h
	l.mu.Unlock()
}

func (l *Logger) Identity() *Identity {
	return l.id
}

func (l *Logger) AddChild(child *Logger) {
	l.mu.Lock()
	if l.children == nil {
		l.children = make(map[uint64][]*Logger, 1)
	}
	// children are group by their identity, i.e a package logger may have many struct logger of same struct because
	// that struct is used in multiple goroutines, we can use string as key, but we use uint64 because we wanted to add
	// hashutil package ... TODO: maybe just change back to string to make people's life easier
	k := child.id.Hash()
	children := l.children[k]
	// avoid putting same pointer twice, though it should never happen if used correctly
	exists := false
	for _, c := range children {
		if c == child {
			exists = true
		}
	}
	if !exists {
		l.children[k] = append(children, child)
	}
	l.mu.Unlock()
}

func (l *Logger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	l.h.HandleLog(PanicLevel, s)
	l.h.Flush()
	panic(s)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Panic(fmt.Sprintf(format, args))
}

func (l *Logger) Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	l.h.HandleLog(FatalLevel, s)
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args))
}

func (l *Logger) PrintTree() {
	l.PrintTreeTo(os.Stdout)
}

// TODO: PrintTree is not implemented
func (l *Logger) PrintTreeTo(w io.Writer) {
	//root := &structure.StringTreeNode{Val: }
}

//// TODO: deal w/ http access log later
//type HttpAccessLogger struct {
//}
