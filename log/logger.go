package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	mu       sync.Mutex
	h        Handler
	level    Level
	fields   Fields
	parent   *Logger
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

func (l *Logger) Debug(args ...interface{}) {
	if l.level >= DebugLevel {
		l.h.HandleLog(DebugLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.level >= InfoLevel {
		l.h.HandleLog(InfoLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) PrintTree() {
	l.PrintTreeTo(os.Stdout)
}

func (l *Logger) PrintTreeTo(w io.Writer) {
	//root := &structure.StringTreeNode{Val: }
}

//// TODO: deal w/ http access log later
//type HttpAccessLogger struct {
//}
