package log

// LoggableStruct is used to inject a logger into the struct,
// the methods for the interface can and should be generated using gommon.
//
// In struct initializer call dlog.NewStructLogger(pkgLogger, structInstancePointer)
type LoggableStruct interface {
	GetLogger() *Logger
	SetLogger(logger *Logger)
	LoggerIdentity(justCallMe func() Identity) Identity
}

func SetLevel(level Level) {
	WalkLogger(func(l *Logger) {
		l.SetLevel(level)
	})
}

func SetHandler(handler Handler) {
	WalkLogger(func(l *Logger) {
		l.SetHandler(handler)
	})
}

func EnableSource() {
	WalkLogger(func(l *Logger) {
		l.EnableSource()
	})
}

func DisableSource() {
	WalkLogger(func(l *Logger) {
		l.DisableSource()
	})
}
