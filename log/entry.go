package log

import (
	"fmt"
	"os"
	"time"
)

type Entry struct {
	Logger  *Logger
	Fields  Fields
	Time    time.Time
	Level   Level
	Message string
}

func (entry *Entry) AddField(key string, value string) {
	entry.Fields[key] = value
}

func (entry *Entry) AddFields(fields Fields) {
	for k, v := range fields {
		entry.Fields[k] = v
	}
}

// This function is not defined with a pointer receiver because we change
// the attribute of struct without using lock, if we use pointer, it would
// become race condition for multiple goroutines.
// see https://github.com/at15/go-learning/issues/3
func (entry Entry) log(level Level, msg string) bool {
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg
	// don't log if it can't pass the filter
	for _, filter := range entry.Logger.Filters[level] {
		if !filter.Filter(&entry) {
			return false
		}
	}
	serialized, err := entry.Logger.Formatter.Format(&entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to serialize, %v\n", err)
		return false
	}
	_, err = entry.Logger.Out.Write(serialized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write, %v\n", err)
		return false
	}
	return true
}

func (entry *Entry) Trace(args ...interface{}) {
	if entry.Logger.Level >= TraceLevel {
		entry.log(TraceLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	// TODO: allow register handlers like logrus
	os.Exit(1)
}

// TODO: maybe Painc should comes after fatal
func (entry *Entry) Painc(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Printf functions
// NOTE: the *f functions does NOT call * functions like logrus does, it just copy and paste

func (entry *Entry) Tracef(format string, args ...interface{}) {
	if entry.Logger.Level >= TraceLevel {
		entry.log(TraceLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
	// TODO: allow register handlers like logrus
	os.Exit(1)
}

// TODO: maybe Painc should comes after fatal
func (entry *Entry) Paincf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
	panic(fmt.Sprint(args...))
}
