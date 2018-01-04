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
	children []*Logger
	id       *Identity
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

func (l *Logger) AddChild(child *Logger) {
	l.mu.Lock()
	l.children = append(l.children, child)
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
