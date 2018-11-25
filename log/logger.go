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
//
// 	logger.Debugf("id is %d", id)
//
// For structural logging (LevelF), Logger passes fields to handlers without any processing.
//
//	logger.DebugF("hi", log.Fields{log.Str("foo", "bar")})
//
// If you want to mix two styles, call fmt.Sprintf before calling DebugF,
//
// 	logger.DebugF(fmt.Sprintf("id is %d", id), log.Fields{log.Str("foo", "bar")})
type Logger struct {
	// mu is a Mutex instead of RWMutex because it's only for avoid concurrent write,
	// for performance reason and the natural of logging, reading stale config is not a big problem,
	// so we don't check mutex on read operation (i.e. log message) and allow race condition
	mu       sync.Mutex
	h        Handler
	level    Level
	source   bool
	children map[string][]*Logger

	// fields contains common context, i.e. the struct is created for a specific task and it has "taskId": 0ac-123
	fields Fields

	id *Identity // use nil so we can have logger without identity
}

// AddField add field to current logger in place, it does NOT make a copy
func (l *Logger) AddField(f Field) *Logger {
	l.mu.Lock()
	// TODO: check dup or not? or may it optional
	l.fields = append(l.fields, f)
	l.mu.Unlock()
	return l
}

// AddFields add fields to current logger in place, it does NOT make a copy
func (l *Logger) AddFields(fields ...Field) *Logger {
	l.mu.Lock()
	// TODO: check dup or not? or may it optional
	l.fields = append(l.fields, fields...)
	l.mu.Unlock()
	return l
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) SetLevel(level Level) *Logger {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
	return l
}

func (l *Logger) SetHandler(h Handler) *Logger {
	l.mu.Lock()
	l.h = h
	l.mu.Unlock()
	return l
}

func (l *Logger) EnableSource() *Logger {
	l.mu.Lock()
	l.source = true
	l.mu.Unlock()
	return l
}

func (l *Logger) DisableSource() *Logger {
	l.mu.Lock()
	l.source = false
	l.mu.Unlock()
	return l
}

// Identity returns the identity set when the logger is created.
// NOTE: caller can modify the identity because all fields are public, but they should NOT do this
func (l *Logger) Identity() Identity {
	if l.id == nil {
		return UnknownIdentity
	}
	return *l.id
}

// Panic calls panic after it writes and flushes the log
func (l *Logger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	if len(l.fields) == 0 {
		l.h.HandleLog(PanicLevel, time.Now(), s, caller(), nil, nil)
	} else {
		l.h.HandleLog(PanicLevel, time.Now(), s, caller(), l.fields, nil)
	}
	l.h.Flush()
	panic(s)
}

// Panicf duplicates instead of calling Panic to keep source line correct
func (l *Logger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if len(l.fields) == 0 {
		l.h.HandleLog(PanicLevel, time.Now(), s, caller(), nil, nil)
	} else {
		l.h.HandleLog(PanicLevel, time.Now(), s, caller(), l.fields, nil)
	}
	l.h.Flush()
	panic(s)
}

// PanicF duplicates instead of calling Panic to keep source line correct
func (l *Logger) PanicF(msg string, fields Fields) {
	if len(l.fields) == 0 {
		l.h.HandleLog(PanicLevel, time.Now(), msg, caller(), nil, fields)
	} else {
		l.h.HandleLog(PanicLevel, time.Now(), msg, caller(), l.fields, fields)
	}
	l.h.Flush()
	panic(msg)
}

// Fatal calls os.Exit(1) after it writes and flushes the log
func (l *Logger) Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	if len(l.fields) == 0 {
		l.h.HandleLog(FatalLevel, time.Now(), s, caller(), nil, nil)
	} else {
		l.h.HandleLog(FatalLevel, time.Now(), s, caller(), l.fields, nil)
	}
	l.h.Flush()
	// TODO: allow user to register hook to do cleanup before exit directly
	os.Exit(1)
}

// Fatalf duplicates instead of calling Fatal to keep source line correct
func (l *Logger) Fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if len(l.fields) == 0 {
		l.h.HandleLog(FatalLevel, time.Now(), s, caller(), nil, nil)
	} else {
		l.h.HandleLog(FatalLevel, time.Now(), s, caller(), l.fields, nil)
	}
	l.h.Flush()
	os.Exit(1)
}

// FatalF duplicates instead of calling Fatal to keep source line correct
func (l *Logger) FatalF(msg string, fields Fields) {
	if len(l.fields) == 0 {
		l.h.HandleLog(FatalLevel, time.Now(), msg, caller(), nil, fields)
	} else {
		l.h.HandleLog(FatalLevel, time.Now(), msg, caller(), l.fields, fields)
	}
	l.h.Flush()
	os.Exit(1)
}

// caller gets source location at runtime, in the future we may generate it at compile time to reduce the
// overhead, though I am not sure what the overhead is without actual benchmark and profiling
// TODO: https://github.com/dyweb/gommon/issues/43
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
