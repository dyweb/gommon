// Package log is not usable yet, see legacy/log
package log

type LoggableStruct interface {
	GetLogger() *Logger
	SetLogger(logger *Logger)
	LoggerIdentity(justCallMe func() *Identity) *Identity
}
