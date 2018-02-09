package log

import "github.com/dyweb/gommon/util/color"

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

func (level Level) String() string {
	return levelStrings[level]
}

func (level Level) ColoredString() string {
	return levelColoredStrings[level]
}

func (level Level) ColoredAlignedUpperString() string {
	return levelColoredAlignedUpperStrings[level]
}
