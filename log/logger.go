package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger is a concrete type instead of interface because most logic is in handler.
// There is NO lock when calling logging methods, handlers may have locks.
// Lock is used when updating logger attributes like Level.
//
// For Printf style logging (Levelf), Logger formats string using fmt.Sprintf before passing it to handlers.
// 	logger.Debugf("id is %d", id)
// For structual logging (LevelF), Logger passes fields to handlers without any processing.
//	logger.DebugF("hi", log.Fields{log.Str("foo", "bar")})
// If you want to mix two styles, call fmt.Sprintf before calling DebugF,
// 	logger.DebugF(fmt.Sprintf("id is %d", id), log.Fields{log.Str("foo", "bar")})
type Logger struct {
	mu    sync.RWMutex
	h     Handler
	level Level
	// TODO: Fields in logger are never used, we are using DebugF to pass temporary fields
	// which does not allow inherit fields from parent logger
	//fields   Fields
	children map[string][]*Logger
	source   bool
	id       *Identity
}

func (l *Logger) Level() Level {
	// TODO: might use the mutex here?
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
	s := fmt.Sprintf(format, args...)
	if !l.source {
		l.h.HandleLog(PanicLevel, time.Now(), s)
	} else {
		l.h.HandleLogWithSource(PanicLevel, time.Now(), s, caller())
	}
	l.h.Flush()
	panic(s)
}

func (l *Logger) PanicF(msg string, fields Fields) {
	if !l.source {
		l.h.HandleLogWithFields(PanicLevel, time.Now(), msg, fields)
	} else {
		l.h.HandleLogWithSourceFields(PanicLevel, time.Now(), msg, caller(), fields)
	}
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

// FIXME: source line is in correct because we call Fatal in Fatalf
func (l *Logger) Fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if !l.source {
		l.h.HandleLog(FatalLevel, time.Now(), s)
	} else {
		l.h.HandleLogWithSource(FatalLevel, time.Now(), s, caller())
	}
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

func (l *Logger) FatalF(msg string, fields Fields) {
	if !l.source {
		l.h.HandleLogWithFields(FatalLevel, time.Now(), msg, fields)
	} else {
		l.h.HandleLogWithSourceFields(FatalLevel, time.Now(), msg, caller(), fields)
	}
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
