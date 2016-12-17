package log

import (
	"io"
	"os"
	"sync"
)

// Level is log level
type Level uint8

const (
	// FatalLevel log error and call `os.Exit(1)`
	FatalLevel Level = iota
	// PanicLevel log error and call painc
	PanicLevel
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

func (level Level) String() string {
	switch level {
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "painc"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	default:
		return "unknown"
	}
}

var AllLevels = []Level{
	FatalLevel,
	PanicLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// Fields is key-value string pair to annotate the log and can be used by filter
type Fields map[string]string

type Logger struct {
	Out       io.Writer
	Formatter Formatter
	Level     Level
	mu        sync.Mutex
	Filters   map[Level]map[string]Filter
}

// NewLogger returns a new logger using StdOut and InfoLevel
func NewLogger() *Logger {
	f := make(map[Level]map[string]Filter, len(AllLevels))
	for _, level := range AllLevels {
		f[level] = make(map[string]Filter, 1)
	}
	l := &Logger{
		Out:       os.Stdout,
		Formatter: &TextFormatter{},
		Level:     InfoLevel,
		Filters:   f,
	}
	return l
}

func (log *Logger) AddFilter(filter Filter, level Level) {
	log.Filters[level][filter.Name()] = filter
}

func (log *Logger) NewEntry() *Entry {
	// TODO: may use pool
	return &Entry{
		Logger: log,
		Fields: make(map[string]string, 1),
	}
}
