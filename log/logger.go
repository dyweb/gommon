package log

import (
	"io"
	"os"
	"sync"

	"strings"

	"github.com/pkg/errors"
)

// Level is log level
type Level uint8

const (
	// FatalLevel log error and call `os.Exit(1)`
	FatalLevel Level = iota
	// PanicLevel log error and call panic
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

// Fields is key-value string pair to annotate the log and can be used by filter
type Fields map[string]string

// Logger is used to set output, formatter and filters, the real log is using Entry
type Logger struct {
	Out            io.Writer
	Formatter      Formatter
	Level          Level
	mu             sync.Mutex // FIXME: this mutex is never used, I guess I was following logrus when I wrote this
	showSourceLine bool
	Filters        map[Level]map[string]Filter
}

// NewLogger returns a new logger using StdOut and InfoLevel
func NewLogger() *Logger {
	f := make(map[Level]map[string]Filter, len(AllLevels))
	for _, level := range AllLevels {
		f[level] = make(map[string]Filter, 1)
	}
	l := &Logger{
		Out:            os.Stdout,
		Formatter:      NewTextFormatter(),
		Level:          InfoLevel,
		Filters:        f,
		showSourceLine: false,
	}
	return l
}

func (log *Logger) ApplyConfig(c *Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if log.Level.String() != c.Level {
		newLevel, err := ParseLevel(c.Level, false)
		if err != nil {
			return errors.Wrapf(err, "can't set logging level to %s", c.Level)
		}
		log.Level = newLevel
	}
	if c.Source {
		log.EnableSourceLine()
	}
	// TODO: set color is not supported in formatter interface, so does time format
	// TODO: pkg filter should also be considered
	return nil
}

// EnableSourceLine add `source` field when logging, it use runtime.Caller(), the overhead has not been measured
func (log *Logger) EnableSourceLine() {
	log.showSourceLine = true
}

// DisableSourceLine does not show `source` field
func (log *Logger) DisableSourceLine() {
	log.showSourceLine = false
}

// AddFilter add a filter to logger, the filter should be simple string check on fields, i.e. PkgFilter check pkg field
func (log *Logger) AddFilter(filter Filter, level Level) {
	log.Filters[level][filter.Name()] = filter
}

// NewEntry returns an Entry with empty Fields
func (log *Logger) NewEntry() *Entry {
	// TODO: may use pool
	return &Entry{
		Logger: log,
		Fields: make(map[string]string, 1),
	}
}

// NewEntryWithPkg returns an Entry with pkg Field set to pkgName, should be used with PkgFilter
func (log *Logger) NewEntryWithPkg(pkgName string) *Entry {
	fields := make(map[string]string, 1)
	fields["pkg"] = pkgName
	return &Entry{
		Logger: log,
		Fields: fields,
	}
}
