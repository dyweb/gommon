package log2

type Handler interface {
	HandleLog() // TODO: we might not use entry, just level, msg, fields instead
	Flush()
}

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
type HandlerFunc func(entry *Entry)

func (f HandlerFunc) HandleLog(e *Entry) {
	f(e)
}