package log

// Level is log level
type Level uint8

const (
	// FatalLevel log error and call `os.Exit(1)`
	// TODO: allow user hook exit?
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

// based on https://github.com/apex/log/blob/master/levels.go
var levelStrings = []string{
	FatalLevel: "fatal",
	PanicLevel: "panic",
	ErrorLevel: "error",
	WarnLevel:  "warn",
	InfoLevel:  "info",
	DebugLevel: "debug",
	TraceLevel: "trace",
}

func (level Level) String() string {
	return levelStrings[level]
}
