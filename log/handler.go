package log

import (
	"fmt"
	"os"
)

type Handler interface {
	HandleLog(level Level, msg string)
	Flush()
}

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
type HandlerFunc func(entry *Entry)

// TODO: why the receiver is value instead of pointer https://github.com/dyweb/gommon/issues/30
func (f HandlerFunc) HandleLog(e *Entry) {
	f(e)
}

type defaultHandler struct {
}

var DefaultHandler Handler = &defaultHandler{}

func (h *defaultHandler) HandleLog(level Level, msg string) {
	// TODO: might use a buffer, since we are just concat string, no special format is needed
	// TODO: time, fields etc.
	fmt.Fprintf(os.Stdout, "%s %s\n", level, msg)
}

func (h *defaultHandler) Flush() {
	// nop
}
