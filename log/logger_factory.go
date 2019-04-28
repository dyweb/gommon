package log

// logger_factory.go creates logger without register them to registry

// TODO: they should all become private ...

func NewPackageLogger() *Logger {
	return newPackageLoggerWithSkip(1)
}

func newPackageLoggerWithSkip(skip int) *Logger {
	id := newIdentityFromCaller(skip + 1)
	return copyOrCreateLogger(nil, &id)
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	id := loggable.LoggerIdentity(func() Identity {
		return newIdentityFromCaller(1)
	})
	l := copyOrCreateLogger(packageLogger, &id)
	loggable.SetLogger(l)
	return l
}

// NewTestLogger does NOT have identity nor handler, it is mainly used for benchmark
func NewTestLogger(level Level) *Logger {
	l := &Logger{
		level: level,
	}
	return l
}

// copyOrCreateLogger inherit handler, level, make copy of fields from parent (if present)
// Or create a new one using default handler, level and no fields
func copyOrCreateLogger(parent *Logger, id *Identity) *Logger {
	child := Logger{
		id: id,
	}
	if parent != nil {
		child.h = parent.h
		child.level = parent.level
		child.source = parent.source
		if len(parent.fields) != 0 {
			child.fields = copyFields(parent.fields)
		}
	} else {
		// TODO: allow customize DefaultHandler
		child.h = DefaultHandler()
		child.level = defaultLevel
	}
	return &child
}
