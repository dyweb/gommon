package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu       sync.RWMutex
	h        Handler
	level    Level
	fields   Fields // TODO: the Fields in logger are never used, we are using DebugF to pass temporary fields
	children map[string][]*Logger
	source   bool
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

func (l *Logger) EnableSource() {
	l.mu.Lock()
	l.source = true
	l.mu.Unlock()
}

func (l *Logger) DisableSource() {
	l.mu.Lock()
	l.source = false
	l.mu.Unlock()
}

func (l *Logger) Identity() *Identity {
	return l.id
}

func (l *Logger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	if !l.source {
		l.h.HandleLog(PanicLevel, time.Now(), s)
	} else {
		l.h.HandleLogWithSource(PanicLevel, time.Now(), s, caller())
	}
	l.h.Flush()
	panic(s)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Panic(fmt.Sprintf(format, args))
}

func (l *Logger) PanicF(msg string, fields Fields) {
	l.h.HandleLogWithFields(PanicLevel, time.Now(), msg, fields)
	l.h.Flush()
	panic(msg)
}

func (l *Logger) Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	if !l.source {
		l.h.HandleLog(FatalLevel, time.Now(), s)
	} else {
		l.h.HandleLogWithSource(FatalLevel, time.Now(), s, caller())
	}
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args))
}

func (l *Logger) FatalF(msg string, fields Fields) {
	l.h.HandleLogWithFields(FatalLevel, time.Now(), msg, fields)
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

func caller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<?>"
		line = 1
	} else {
		last := strings.LastIndex(file, "/")
		file = file[last+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
