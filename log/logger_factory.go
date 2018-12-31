package log

// logger_factory.go creates logger without register them to registry

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

// Deprecated: use Copy method on package logger
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

// Deprecated: use Copy method on struct logger
func NewMethodLogger(structLogger *Logger) *Logger {
	id := NewIdentityFromCaller(1)
	l := &Logger{
		id: &id,
	}
	return newLogger(structLogger, l)
}

// NewTestLogger does not have identity and handler, it is mainly used for benchmark test
func NewTestLogger(level Level) *Logger {
	l := &Logger{
		level: level,
	}
	return l
}

func newLogger(parent *Logger, child *Logger) *Logger {
	if parent != nil {
		child.h = parent.h
		child.level = parent.level
		child.source = parent.source
		if len(parent.fields) != 0 {
			fields := make([]Field, len(parent.fields))
			copy(fields, parent.fields)
			child.fields = fields
		}
	} else {
		child.h = DefaultHandler()
		child.level = InfoLevel
	}
	return child
}
