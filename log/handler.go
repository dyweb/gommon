package log

import (
	"fmt"
	"os"
	"sync"
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
	// TODO: it seems in go, both os.Stderr and os.Stdout are not (line) buffered
	fmt.Fprintf(os.Stderr, "%s %s\n", level, msg)
}

func (h *defaultHandler) Flush() {
	// TODO: don't know if is needed, will there be any different if stderr is redirected to a file
	os.Stderr.Sync()
}

// unlike log v1 entry is only used for test, it is not passed around
type entry struct {
	level Level
	msg   string
}

type testHandler struct {
	mu      sync.RWMutex
	entries []entry
}

var TestHandler Handler = &testHandler{}

func NewTestHandler() *testHandler {
	return &testHandler{}
}

func (h *testHandler) HandleLog(level Level, msg string) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, msg: msg})
	h.mu.Unlock()
}

func (h *testHandler) Flush() {
	// nop
}

func (h *testHandler) HasLog(level Level, msg string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, e := range h.entries {
		if e.level == level && e.msg == msg {
			return true
		}
	}
	return false
}
