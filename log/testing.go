package log

// test.go contains helpers for both internal and external testing

import (
	"sync"
	"time"
)

// entry is only used for test, we do not use it as contract for interface
type entry struct {
	level   Level
	time    time.Time
	msg     string
	context Fields
	fields  Fields
	source  Caller
}

var _ Handler = (*TestHandler)(nil)

// TestHandler stores log as entry, its slice is protected by a RWMutex and safe for concurrent use
type TestHandler struct {
	mu      sync.RWMutex
	entries []entry
}

// NewTestHandler returns a test handler, it should only be used in test,
// a concrete type instead of Handler interface is returned to reduce unnecessary type cast in test
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) HandleLog(level Level, time time.Time, msg string, source Caller, context Fields, fields Fields) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg, source: source, context: CopyFields(context), fields: CopyFields(fields)})
	h.mu.Unlock()
}

// Flush implements Handler interface
func (h *TestHandler) Flush() {
	// nop
}

// HasLog checks if a log with specified level and message exists in slice
// TODO: support field, source etc.
func (h *TestHandler) HasLog(level Level, msg string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, e := range h.entries {
		if e.level == level && e.msg == msg {
			return true
		}
	}
	return false
}

func (h *TestHandler) getLogByMessage(msg string) (entry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, e := range h.entries {
		if e.msg == msg {
			return e, true
		}
	}
	return entry{}, false
}
