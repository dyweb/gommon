package log2

import "sync"

type Logger interface {
	Children() []Logger
	SetHandler(h Handler)
}

type BaseLogger struct {
	mu sync.Mutex
	h  Handler
}

type AppLogger struct {
}

type PkgLogger struct {
}

type FuncLogger struct {
	Parent *PkgLogger
}

type StructLogger struct {
	Parent *PkgLogger
}

type MethodLogger struct {
	Parent *StructLogger
}

type HttpAccessLogger struct {
}

func (l *BaseLogger) SetHandler(h Handler) {
	l.mu.Lock()
	l.h = h
	l.mu.Unlock()
}
