package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Entry is the real logger
type Entry struct {
	Logger     *Logger
	Pkg        string
	EntryLevel Level
	Fields     Fields
	Time       time.Time
	Level      Level
	Message    string
}

// SetPkgAlias allows use shorter name for pkg when logging
func (entry *Entry) SetPkgAlias(alias string) {
	entry.Fields["pkg"] = alias
}

func (entry *Entry) SetEntryLevel(s string) error {
	newLevel, err := ParseLevel(s, false)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("can't set logging level to %s", s))
	}
	entry.EntryLevel = newLevel
	return nil
}

// AddField adds tag to entry
func (entry *Entry) AddField(key string, value string) {
	entry.Fields[key] = value
}

// AddFields adds multiple tags to entry
func (entry *Entry) AddFields(fields Fields) {
	for k, v := range fields {
		entry.Fields[k] = v
	}
}

// DeleteField remove a tag from entry, this was added for benchmark to remove the automatically added pkg tag when using RegisterPkg
func (entry *Entry) DeleteField(key string) {
	delete(entry.Fields, key)
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
		if !filter.Accept(&entry) {
			return false
		}
	}
	// add source code line if required
	if entry.Logger.showSourceLine {
		// TODO: what if the user also have tag called source
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "<?>"
			line = 1
		} else {
			lastSlash := strings.LastIndex(file, "/")
			file = file[lastSlash+1:]
		}
		entry.AddField("source", fmt.Sprintf("%s:%d", file, line))
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

func (entry *Entry) Panic(args ...interface{}) {
	if entry.EntryLevel >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.EntryLevel >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	// TODO: allow register handlers like logrus
	os.Exit(1)
}

// Printf functions
// NOTE: the *f functions does NOT call * functions like logrus does, it just copy and paste

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.EntryLevel >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
	panic(fmt.Sprint(args...))
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.EntryLevel >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprintf(format, args...))
	}
	// TODO: allow register handlers like logrus
	os.Exit(1)
}
