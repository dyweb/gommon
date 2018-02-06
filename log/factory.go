package log

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
	l := &Logger{
		id: NewIdentityFromCaller(skip + 1),
	}
	return newLogger(nil, l)
}

func NewFunctionLogger(packageLogger *Logger) *Logger {
	l := &Logger{
		id: NewIdentityFromCaller(1),
	}
	return newLogger(packageLogger, l)
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	l := &Logger{
		id: loggable.LoggerIdentity(func() *Identity {
			return NewIdentityFromCaller(1)
		}),
	}
	l = newLogger(packageLogger, l)
	loggable.SetLogger(l)
	return l
}

func NewMethodLogger(structLogger *Logger) *Logger {
	l := &Logger{
		id: NewIdentityFromCaller(1),
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
