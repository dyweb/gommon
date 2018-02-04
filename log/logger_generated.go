// Code generated by gommon from log/logger_generated.go.tmpl DO NOT EDIT.
package log

import (
	"fmt"
)

func (l *Logger) IsTraceEnabled() bool {
	return l.level >= TraceLevel
}

func (l *Logger) Trace(args ...interface{}) {
	if l.level >= TraceLevel {
		l.h.HandleLog(TraceLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.level >= TraceLevel {
		l.h.HandleLog(TraceLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) IsDebugEnabled() bool {
	return l.level >= DebugLevel
}

func (l *Logger) Debug(args ...interface{}) {
	if l.level >= DebugLevel {
		l.h.HandleLog(DebugLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level >= DebugLevel {
		l.h.HandleLog(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) IsInfoEnabled() bool {
	return l.level >= InfoLevel
}

func (l *Logger) Info(args ...interface{}) {
	if l.level >= InfoLevel {
		l.h.HandleLog(InfoLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level >= InfoLevel {
		l.h.HandleLog(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) IsWarnEnabled() bool {
	return l.level >= WarnLevel
}

func (l *Logger) Warn(args ...interface{}) {
	if l.level >= WarnLevel {
		l.h.HandleLog(WarnLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level >= WarnLevel {
		l.h.HandleLog(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) IsErrorEnabled() bool {
	return l.level >= ErrorLevel
}

func (l *Logger) Error(args ...interface{}) {
	if l.level >= ErrorLevel {
		l.h.HandleLog(ErrorLevel, fmt.Sprint(args...))
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		l.h.HandleLog(ErrorLevel, fmt.Sprintf(format, args...))
	}
}
