package log2

import (
	"fmt"
	"runtime"
	"github.com/dyweb/gommon/util/runtimeutil"
)

type LoggerType uint8

const (
	UnknownLogger     LoggerType = iota
	ApplicationLogger
	PackageLogger
	FunctionLogger
	StructLogger
	MethodLogger
)

var loggerTypeStrings = []string{
	UnknownLogger:     "unk",
	ApplicationLogger: "app",
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

const MagicStructLoggerMethod = "LoggerIdentity"

type LoggableStruct interface {
	LoggerIdentity(justCallMe func() *Identity) *Identity
}

func NewPackageLogger() *Logger {
	// TODO: the identity should work just fine
	return nil
}

func NewFunctionLogger(packageLogger *Logger) *Logger {
	// TODO: set the parent, runtime can handle the identity
	return nil
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	// TODO: pass a function to loggable
	return nil
}

func NewMethodLogger(structLogger *Logger) *Logger {
	// TODO: set the parent, runtime can handle the identity
	return nil
}

// see https://github.com/dyweb/gommon/issues/32
// based on https://github.com/go-stack/stack/blob/master/stack.go#L29:51
// TODO: not sure if calling two Next without checking the more value works for other go version
func NewIdentityFromCaller(skip int) *Identity {
	var pcs [3]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	f, _ := frames.Next()
	f, _ = frames.Next()

	fmt.Println("Function", f.Function)
	fmt.Println("File", f.File)
	fmt.Println("Line", f.Line)
	fmt.Println("Func.Name()", f.Func.Name())

	return nil
	//return &Identity{
	//	Function: f.Function,
	//	File:     f.File,
	//	Line:     f.Line,
	//}
}

func NewIdentityFromCaller2() *Identity {
	// TODO: handle package level call where there is no function
	frame := runtimeutil.GetCallerFrame(1)
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
		if function == MagicStructLoggerMethod {
			tpe = StructLogger
		}
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
