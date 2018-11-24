package log

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultTimeStampFormat = time.RFC3339
)

// Handler formats log message and writes to underlying storage, stdout, file, remote server etc.
// It MUST be thread safe because logger calls handler concurrently without any locking.
// There is NO log entry struct in gommon/log, which is used in many logging packages, the reason is
// if extra field is added to the interface, compiler would throw error on stale handler implementations.
type Handler interface {
	// HandleLog accepts level, log time, formatted log message
	HandleLog(level Level, time time.Time, msg string)
	// HandleLogWithSource accepts formatted source line of log i.e., http.go:13
	// TODO: pass frame instead of string so handler can use trace for error handling?
	HandleLogWithSource(level Level, time time.Time, msg string, source string)
	// HandleLogWithFields accepts fields with type hint,
	// implementation should inspect the type field instead of using reflection
	HandleLogWithFields(level Level, time time.Time, msg string, fields Fields)
	// HandleLogWithSourceFields accepts both source and fields
	HandleLogWithSourceFields(level Level, time time.Time, msg string, source string, fields Fields)
	// Flush writes the buffer to underlying storage
	Flush()
}

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
//type HandlerFunc func(level Level, msg string)

// TODO: why the receiver is value instead of pointer https://github.com/dyweb/gommon/issues/30
//func (f HandlerFunc) HandleLog(level Level, msg string) {
//	f(level, msg)
//}

var _ Syncer = (*os.File)(nil)

// Syncer is implemented by os.File, handler implementation should check this interface and call Sync
// if they support using file as sink
type Syncer interface {
	Sync() error
}

//// TODO: handler for http access log, this should be in extra package
//type HttpAccessLogger struct {
//}`

var defaultHandler = NewIOHandler(os.Stderr)

// DefaultHandler returns the singleton defaultHandler instance, which logs to stdout in text format
func DefaultHandler() Handler {
	return defaultHandler
}

// IOHandler writes log to io.Writer, default handler uses os.Stderr
type IOHandler struct {
	w io.Writer
}

// NewIOHandler
func NewIOHandler(w io.Writer) Handler {
	return &IOHandler{w: w}
}

// TODO: performance (which is not a major concern now ...)
// - when using raw byte slice, have a correct length, fields can also return length required
// - is calling level.String() faster than %s level
// - use buffer (pool)
// TODO: correctness
// - in go, both os.Stderr and os.Stdout are not (line) buffered
// - what would happen if os.Stderr.Close()
// - os.Stderr.Sync() will there be any different if stderr/stdout is redirected to a file

// HandleLog implements Handler interface
func (h *IOHandler) HandleLog(level Level, time time.Time, msg string) {
	b := formatHead(level, time, msg)
	b = append(b, '\n')
	h.w.Write(b)
}

// HandleLogWithSource implements Handler interface
func (h *IOHandler) HandleLogWithSource(level Level, time time.Time, msg string, source string) {
	b := formatHeadWithSource(level, time, msg, source)
	b = append(b, '\n')
	h.w.Write(b)
}

// HandleLogWithFields implements Handler interface
func (h *IOHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	// we use raw slice instead of bytes buffer because we need to use strconv.Append*, which requires raw slice
	b := formatHead(level, time, msg)
	b = append(b, ' ')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

// HandleLogWithSourceFields implements Handler interface
func (h *IOHandler) HandleLogWithSourceFields(level Level, time time.Time, msg string, source string, fields Fields) {
	b := formatHeadWithSource(level, time, msg, source)
	b = append(b, ' ')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

// Flush implements Handler interface
func (h *IOHandler) Flush() {
	if s, ok := h.w.(Syncer); ok {
		s.Sync()
	}
}

// entry is only used for test, it is not passed around like other loging packages
type entry struct {
	level  Level
	time   time.Time
	msg    string
	fields Fields
	source string
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

// HandleLog implements Handler interface
func (h *TestHandler) HandleLog(level Level, time time.Time, msg string) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg})
	h.mu.Unlock()
}

// HandleLogWithSource implements Handler interface
func (h *TestHandler) HandleLogWithSource(level Level, time time.Time, msg string, source string) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg, source: source})
	h.mu.Unlock()
}

// HandleLogWithFields implements Handler interface
func (h *TestHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg, fields: fields})
	h.mu.Unlock()
}

// HandleLogWithSourceFields implements Handler interface
func (h *TestHandler) HandleLogWithSourceFields(level Level, time time.Time, msg string, source string, fields Fields) {
	h.mu.Lock()
	h.entries = append(h.entries, entry{level: level, time: time, msg: msg, source: source, fields: fields})
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

// no need to use fmt.Printf since we don't need any format
func formatHead(level Level, time time.Time, msg string) []byte {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, msg...)
	return b
}

// we have a new function because source sits between time and msg in output, instead of after msg
// i.e. info 2018-02-04T21:03:20-08:00 main.go:18 show me the line
func formatHeadWithSource(level Level, time time.Time, msg string, source string) []byte {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg)+len(source))
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, source...)
	b = append(b, ' ')
	b = append(b, msg...)
	return b
}

// it has an extra tailing space, which can be updated inplace to a \n
func formatFields(b []byte, fields Fields) []byte {
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
	return b
}
