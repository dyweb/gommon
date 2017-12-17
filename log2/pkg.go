// Package log2 should be renamed to log once the migration is finished
package log2

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
type HandlerFunc func(entry *Entry)

func (f HandlerFunc) HandleLog(e *Entry) {
	f(e)
}