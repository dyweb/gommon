// +build ignore

package log

import (
	"strings"

	"github.com/pkg/errors"
)

// Level is log level
type Level uint8

const (
	// FatalLevel log error and call `os.Exit(1)`
	FatalLevel Level = iota
	// PanicLevel log error and call `panic`
	PanicLevel
	// ErrorLevel log error
	ErrorLevel
	// WarnLevel log warning
	WarnLevel
	// InfoLevel log info
	InfoLevel
	// DebugLevel log debug message, user should enable DebugLevel logging when report bug
	DebugLevel
	// TraceLevel is very verbose, user should enable it only on packages they are currently investing instead of globally
	TraceLevel
)

// ShortUpperString returns the first 4 characters of a level in upper case
func (level Level) ShortUpperString() string {
	switch level {
	case FatalLevel:
		return "FATA"
	case PanicLevel:
		return "PANI"
	case ErrorLevel:
		return "ERRO"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBU"
	case TraceLevel:
		return "TRAC"
	default:
		return "UNKN"
	}
}

func (level Level) String() string {
	switch level {
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
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

// ParseLevel match the level string with Level, it will use strings.HasPrefix in non strict mode
func ParseLevel(s string, strict bool) (Level, error) {
	str := strings.ToLower(s)
	if strict {
		switch str {
		case "fatal":
			return FatalLevel, nil
		case "panic":
			return PanicLevel, nil
		case "error":
			return ErrorLevel, nil
		case "warn":
			return WarnLevel, nil
		case "info":
			return InfoLevel, nil
		case "debug":
			return DebugLevel, nil
		case "trace":
			return TraceLevel, nil
		default:
			return Level(250), errors.Errorf("unknown log level %s", str)
		}
	}
	switch {
	case strings.HasPrefix(str, "f"):
		return FatalLevel, nil
	case strings.HasPrefix(str, "p"):
		return PanicLevel, nil
	case strings.HasPrefix(str, "e"):
		return ErrorLevel, nil
	case strings.HasPrefix(str, "w"):
		return WarnLevel, nil
	case strings.HasPrefix(str, "i"):
		return InfoLevel, nil
	case strings.HasPrefix(str, "d"):
		return DebugLevel, nil
	case strings.HasPrefix(str, "t"):
		return TraceLevel, nil
	default:
		return Level(250), errors.Errorf("unknown log level %s", str)
	}
}

// AllLevels includes all the logging level
var AllLevels = []Level{
	FatalLevel,
	PanicLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}
