// Package log is a structured logging with fine grained control
package log

//// TODO: deal w/ http access log later
//type HttpAccessLogger struct {
//}

type LoggableStruct interface {
	GetLogger() *Logger
	SetLogger(logger *Logger)
	LoggerIdentity(justCallMe func() *Identity) *Identity
}
