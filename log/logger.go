package log

import (
	"io"
	"os"
	"sync"
)

// Level is log level
type Level uint8

const (
	// PanicLevel log error and call painc
	PanicLevel Level = iota
	// FatalLevel log error and call `os.Exit(1)`
	FatalLevel
	// ErrorLevel log error
	ErrorLevel
	// WarnLevel log warning
	WarnLevel
	// InfoLevel log info
	InfoLevel
	// DebugLevel log debug message, user should enable DebugLevel logging when report bug
	DebugLevel
	// TraceLevel is used by developer only, user should stop at DebugLevel if they just want to report
	TraceLevel
)

// Fields is key-value string pair to annotate the log and can be used by filter
type Fields map[string]string

type Logger struct {
	Out   io.Writer
	Level Level
	mu    sync.Mutex
}

// NewLogger returns a new logger using StdOut and InfoLevel
func NewLogger() *Logger {
	return &Logger{
		Out:   os.Stdout,
		Level: InfoLevel,
	}
}

func (log *Logger) AddFilter(filter Filter, level Level) {
}
