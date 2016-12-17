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
