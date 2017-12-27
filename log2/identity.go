package log2

import (
	"github.com/dyweb/gommon/util/runtimeutil"
)

type LoggerType uint8

const (
	UnknownLogger LoggerType = iota
	ApplicationLogger
	PackageLogger
	FunctionLogger
	StructLogger
	MethodLogger
)

var loggerTypeStrings = []string{
	UnknownLogger:     "unk",
	ApplicationLogger: "app",
	PackageLogger:     "pkg",
	FunctionLogger:    "func",
	StructLogger:      "struct",
	MethodLogger:      "method",
}

func (tpe LoggerType) String() string {
	return loggerTypeStrings[tpe]
}

// Identity is based where the logger is initialized, it is NOT exactly where the log happens.
// It is used for applying filter rules and print logger hierarchy.
// TODO: example
type Identity struct {
	Package  string
	Function string
	Struct   string
	File     string
	Line     int
	Type     LoggerType
}

var UnknownIdentity = Identity{Package: "unk", Type: UnknownLogger}

const MagicStructLoggerFunctionName = "LoggerIdentity"
const MagicPackageLoggerFunctionName = "init"

type LoggableStruct interface {
	LoggerIdentity(justCallMe func() *Identity) *Identity
}

func NewPackageLogger() *Logger {
	return &Logger{
		id: NewIdentityFromCaller(1),
		h:  DefaultHandler,
	}
}

func NewFunctionLogger(packageLogger *Logger) *Logger {
	l := &Logger{
		parent: packageLogger,
		id:     NewIdentityFromCaller(1),
	}
	return newLogger(packageLogger, l)
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	l := &Logger{
		parent: packageLogger,
		id: loggable.LoggerIdentity(func() *Identity {
			return NewIdentityFromCaller(1)
		}),
	}
	return newLogger(packageLogger, l)
}

func NewMethodLogger(structLogger *Logger) *Logger {
	l := &Logger{
		parent: structLogger,
		id:     NewIdentityFromCaller(1),
	}
	return newLogger(structLogger, l)
}

func newLogger(parent *Logger, child *Logger) *Logger {
	if parent != nil {
		// TODO: might have a method called add children on Logger
		parent.children = append(parent.children, child)
		child.h = parent.h
	} else {
		child.h = DefaultHandler
	}
	return child
}

// TODO: document all the black magic here ...
// https://github.com/dyweb/gommon/issues/32
func NewIdentityFromCaller(skip int) *Identity {
	frame := runtimeutil.GetCallerFrame(skip + 1)
	var (
		pkg      string
		function string
		st       string
	)
	tpe := UnknownLogger
	pkg, function = runtimeutil.SplitPackageFunc(frame.Function)
	tpe = FunctionLogger
	// NOTE: we distinguish a struct logger and method logger using the magic name,
	// which is normally the case as long as you are using the factory methods to create logger
	// otherwise you have to manually update the type of logger
	if runtimeutil.IsMethod(function) {
		st, function = runtimeutil.SplitStructMethod(function)
		tpe = MethodLogger
		if function == MagicStructLoggerFunctionName {
			tpe = StructLogger
		}
	} else if MagicPackageLoggerFunctionName == function {
		tpe = PackageLogger
	}

	return &Identity{
		Package:  pkg,
		Function: function,
		Struct:   st,
		File:     frame.File,
		Line:     frame.Line,
		Type:     tpe,
	}
}

func (id *Identity) Diff(parent *Identity) string {
	if parent == nil {
		return id.Package
	}
	// TODO: might return full package
	if id.Package != parent.Package {
		return id.Package
	}
	// TODO: this may not be the desired behaviour
	return id.File
}
