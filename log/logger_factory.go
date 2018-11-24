package log

// used for benchmark test
func NewTestLogger(level Level) *Logger {
	l := &Logger{
		level: level,
	}
	return l
}

func NewApplicationLogger() *Logger {
	l := NewPackageLoggerWithSkip(1)
	l.id.Type = ApplicationLogger
	return l
}

func NewLibraryLogger() *Logger {
	l := NewPackageLoggerWithSkip(1)
	l.id.Type = LibraryLogger
	return l
}

func NewPackageLogger() *Logger {
	return NewPackageLoggerWithSkip(1)
}

func NewPackageLoggerWithSkip(skip int) *Logger {
	id := NewIdentityFromCaller(skip + 1)
	l := &Logger{
		id: &id,
	}
	return newLogger(nil, l)
}

func NewFunctionLogger(packageLogger *Logger) *Logger {
	id := NewIdentityFromCaller(1)
	l := &Logger{
		id: &id,
	}
	return newLogger(packageLogger, l)
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	id := loggable.LoggerIdentity(func() Identity {
		return NewIdentityFromCaller(1)
	})
	l := &Logger{
		id: &id,
	}
	l = newLogger(packageLogger, l)
	loggable.SetLogger(l)
	return l
}

func NewMethodLogger(structLogger *Logger) *Logger {
	id := NewIdentityFromCaller(1)
	l := &Logger{
		id: &id,
	}
	return newLogger(structLogger, l)
}

func newLogger(parent *Logger, child *Logger) *Logger {
	if parent != nil {
		parent.AddChild(child)
		child.h = parent.h
		child.level = parent.level
		child.source = parent.source
	} else {
		child.h = DefaultHandler()
		child.level = InfoLevel
	}
	return child
}
