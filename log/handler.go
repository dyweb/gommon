package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	defaultTimeStampFormat = time.RFC3339
)

type Handler interface {
	// TODO: pass pointer for fields?
	HandleLog(level Level, time time.Time, msg string)
	//HandleLogWithSource(source string, level Level, time time.Time, msg string)
	HandleLogWithFields(level Level, time time.Time, msg string, fields Fields)
	//HandleLogWithSourceFields(source string, level Level, time time.Time, msg string, fields Fields)
	Flush()
}

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
//type HandlerFunc func(level Level, msg string)

// TODO: why the receiver is value instead of pointer https://github.com/dyweb/gommon/issues/30
//func (f HandlerFunc) HandleLog(level Level, msg string) {
//	f(level, msg)
//}

type defaultHandler struct {
}

var DefaultHandler Handler = &defaultHandler{}

func (h *defaultHandler) HandleLog(level Level, time time.Time, msg string) {
	// TODO: might use a buffer, since we are just concat string, no special format is needed
	// TODO: fields etc.
	// TODO: is calling level.String() faster than %s level
	// TODO: it seems in go, both os.Stderr and os.Stdout are not (line) buffered
	fmt.Fprintf(os.Stderr, "%s %s %s\n", level.String(), time.Format(defaultTimeStampFormat), msg)
}

func (h *defaultHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	// FIXME: deal with fields
	fmt.Fprintf(os.Stderr, "%s %s %s\n", level.String(), time.Format(defaultTimeStampFormat), msg)
}

func (h *defaultHandler) Flush() {
	// TODO: don't know if is needed, will there be any different if stderr/stdout is redirected to a file
	os.Stderr.Sync()
}

// unlike log v1 entry is only used for test, it is not passed around
type entry struct {
	level Level
	time  time.Time
	msg   string
	// TODO: fields
}

type testHandler struct {
	mu      sync.RWMutex
	entries []entry
}

func NewTestHandler() *testHandler {
	return &testHandler{}
}

func (h *testHandler) HandleLog(level Level, time time.Time, msg string) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg})
	h.mu.Unlock()
}

func (h *testHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	h.mu.Lock()
	// FIXME: deal with fields
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg})
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
