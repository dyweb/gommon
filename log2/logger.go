package log2

import (
	"sync"
	"fmt"
)

//type Logger interface {
//	SetHandler(h Handler)
//}

//type NonTerminalLogger interface {
//	Children() []Logger
//}

type Logger struct {
	mu       sync.Mutex
	h        Handler
	level    Level
	fields   Fields
	parent   *Logger
	children []*Logger
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

//type AppLogger struct {
//	*BaseLogger
//}
//
//type PkgLogger struct {
//	*BaseLogger
//}
//
//type FuncLogger struct {
//	*BaseLogger
//	Parent *PkgLogger
//}
//
//type StructLogger struct {
//	*BaseLogger
//	Parent   *PkgLogger
//	children []*MethodLogger
//}
//
//type MethodLogger struct {
//	*BaseLogger
//	Parent *StructLogger
//}
//
//// TODO: deal w/ http access log later
//type HttpAccessLogger struct {
//}

//func (l *BaseLogger) SetHandler(h Handler) {
//	l.mu.Lock()
//	l.h = h
//	l.mu.Unlock()
//}
