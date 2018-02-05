package log

import (
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultTimeStampFormat = time.RFC3339
)

type Handler interface {
	HandleLog(level Level, time time.Time, msg string)
	//HandleLogWithSource(source string, level Level, time time.Time, msg string)
	// TODO: pass pointer for fields?
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

type stderrHandler struct {
}

var defaultHandler = &stderrHandler{}

func DefaultHandler() Handler {
	return defaultHandler
}

// TODO: performance (which is not a major concern now ...)
// - when using raw byte slice, have a correct length, fields can also return length required
// - is calling level.String() faster than %s level
// - use buffer (pool)
// TODO: correctness
// - in go, both os.Stderr and os.Stdout are not (line) buffered

func (h *stderrHandler) HandleLog(level Level, time time.Time, msg string) {
	// no need to use fmt.Printf since we don't need any format
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, msg...)
	b = append(b, '\n')
	os.Stderr.Write(b)
}

func (h *stderrHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	// we use raw slice instead of bytes buffer because we need to use strconv.Append*, which requires raw slice
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, msg...)
	for _, f := range fields {
		b = append(b, f.Key...)
		b = append(b, '=')
		switch f.Type {
		case IntType:
			b = strconv.AppendInt(b, f.Int, 10)
		case StringType:
			b = append(b, f.Str...)
		}
		b = append(b, ' ')
	}
	b[len(b)-1] = '\n'
	os.Stderr.Write(b)
}

func (h *stderrHandler) Flush() {
	// TODO: don't know if is needed, will there be any different if stderr/stdout is redirected to a file
	os.Stderr.Sync()
}

// unlike log v1 entry is only used for test, it is not passed around
type entry struct {
	level  Level
	time   time.Time
	msg    string
	fields Fields
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
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg, fields: fields})
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
