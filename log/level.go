package log

import "github.com/dyweb/gommon/util/color"

// Level is log level
// TODO: allow change default logging level at compile time
type Level uint8

const (
	// FatalLevel log error and call `os.Exit(1)`
	// TODO: allow user to add hooks before calling os.Exit?
	FatalLevel Level = iota
	// PanicLevel log error and call `panic`
	PanicLevel
	// ErrorLevel log error and do nothing
	// TODO: add integration with errors package
	ErrorLevel
	// WarnLevel log warning that is often ignored
	WarnLevel
	// InfoLevel log info
	InfoLevel
	// DebugLevel log debug message, user should enable DebugLevel logging when report bug
	DebugLevel
	// TraceLevel is very verbose, user should enable it only on packages they are currently investing instead of globally
	// TODO: add compile flag to use empty trace logger implementation to eliminate the call at runtime
	TraceLevel
)

// PrintLevel is for library/application that requires a Printf based logger interface
const PrintLevel = InfoLevel

var levelColoredStrings = []string{
	FatalLevel: color.RedStart + "fatal" + color.End,
	PanicLevel: color.RedStart + "panic" + color.End,
	ErrorLevel: color.RedStart + "error" + color.End,
	WarnLevel:  color.YellowStart + "warn" + color.End,
	InfoLevel:  color.BlueStart + "info" + color.End,
	DebugLevel: color.GrayStart + "debug" + color.End,
	TraceLevel: color.GrayStart + "trace" + color.End,
}

var levelColoredAlignedUpperStrings = []string{
	FatalLevel: color.RedStart + "FATA" + color.End,
	PanicLevel: color.RedStart + "PANI" + color.End,
	ErrorLevel: color.RedStart + "ERRO" + color.End,
	WarnLevel:  color.YellowStart + "WARN" + color.End,
	InfoLevel:  color.BlueStart + "INFO" + color.End,
	DebugLevel: color.GrayStart + "DEBU" + color.End,
	TraceLevel: color.GrayStart + "TRAC" + color.End,
}

// String returns log level in lower case and not aligned in length, i.e. fatal, warn
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

// AlignedUpperString returns log level with fixed length of 4 in uppercase, i.e. FATA, WARN
func (level Level) AlignedUpperString() string {
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
		return "DEBG" // TODO: or DEBU
	case TraceLevel:
		return "TRAC"
	default:
		return "UNKN"
	}
}

// TODO: use switch and generate the function ... or just generate it manually
// ColoredString returns level string wrapped by terminal color characters, only works on *nix
func (level Level) ColoredString() string {
	return levelColoredStrings[level]
}

// ColoredAlignedUpperString returns fixed length level string in uppercase,
// wrapped by terminal color characters, only works on *nix
func (level Level) ColoredAlignedUpperString() string {
	return levelColoredAlignedUpperStrings[level]
}
