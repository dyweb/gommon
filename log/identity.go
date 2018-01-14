package log

import (
	"fmt"

	"github.com/dyweb/gommon/util/hashutil"
	"github.com/dyweb/gommon/util/runtimeutil"
)

type LoggerType uint8

const (
	UnknownLogger LoggerType = iota
	ApplicationLogger
	LibraryLogger
	PackageLogger
	FunctionLogger
	StructLogger
	MethodLogger
)

var loggerTypeStrings = []string{
	UnknownLogger:     "unk",
	ApplicationLogger: "app",
	LibraryLogger:     "lib",
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
	} else if function == MagicPackageLoggerFunctionName {
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
