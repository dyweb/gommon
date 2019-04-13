package log

import (
	"fmt"

	"github.com/dyweb/gommon/util/hashutil"
	"github.com/dyweb/gommon/util/runtimeutil"
)

// LoggerType can be used for filtering loggers, it is set when creating logger
type LoggerType uint8

const (
	UnknownLogger LoggerType = iota
	// PackageLogger is normally singleton in entire package
	// We used to have application and library logger but they are replaced by registry
	PackageLogger
	FunctionLogger
	StructLogger
	MethodLogger
)

var loggerTypeStrings = []string{
	UnknownLogger:  "unk",
	PackageLogger:  "pkg",
	FunctionLogger: "func",
	StructLogger:   "struct",
	MethodLogger:   "method",
}

func (tpe LoggerType) String() string {
	return loggerTypeStrings[tpe]
}

// Identity is set based on logger's initialization location,
// it is close to, but NOT exactly same as location of actual log.
// It is used for applying filter rules and print logger hierarchy.
type Identity struct {
	Package  string
	Function string
	Struct   string
	File     string
	Line     int
	Type     LoggerType
}

var UnknownIdentity = Identity{Package: "unk", Type: UnknownLogger}

const (
	MagicStructLoggerFunctionName  = "LoggerIdentity"
	MagicPackageLoggerFunctionName = "init"
	// a hack for new init func name after 1.12
	// See https://github.com/dyweb/gommon/issues/108
	MagicPackageLoggerFunctionNameGo112 = "init.ializers"
)

// TODO: document all the black magic here ...
// https://github.com/dyweb/gommon/issues/32
func NewIdentityFromCaller(skip int) Identity {
	frame := runtimeutil.GetCallerFrame(skip + 1)
	var (
		pkg      string
		function string
		st       string
	)
	tpe := UnknownLogger
	// TODO: does it handle vendor correctly, and what about vgo ...
	pkg, function = runtimeutil.SplitPackageFunc(frame.Function)
	tpe = FunctionLogger
	if function == MagicPackageLoggerFunctionNameGo112 || function == MagicPackageLoggerFunctionName {
		// https://github.com/dyweb/gommon/issues/108 there are two names, init and init.ializers (after go1.12)
		tpe = PackageLogger
	} else if runtimeutil.IsMethod(function) {
		// NOTE: we distinguish a struct logger and method logger using the magic name,
		// which is normally the case as long as you are using the factory methods to create logger
		// otherwise you have to manually update the type of logger
		st, function = runtimeutil.SplitStructMethod(function)
		tpe = MethodLogger
		if function == MagicStructLoggerFunctionName {
			tpe = StructLogger
		}
	}

	return Identity{
		Package:  pkg,
		Function: function,
		Struct:   st,
		File:     frame.File,
		Line:     frame.Line,
		Type:     tpe,
	}
}

func (id *Identity) Hash() uint64 {
	// assume no one create two logger in one line and no non ascii characters in file path
	return hashutil.HashStringFnv64a(id.SourceLocation())
}

func (id *Identity) SourceLocation() string {
	return fmt.Sprintf("%s:%d", id.File, id.Line)
}

func (id *Identity) String() string {
	return fmt.Sprintf("%s logger %s:%d", id.Type, id.File, id.Line)
}

// TODO: this is used for print tree like structure ... it's hard to maintain exact parent and child logger due to cycle import
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
