package log

import (
	"io"
	"os"
	"strconv"
	"time"
)

// handler.go contains Handler interface and builtin implementations

const (
	defaultTimeStampFormat = time.RFC3339
)

// Handler formats log message and writes to underlying storage, stdout, file, remote server etc.
// It MUST be thread safe because logger calls handler concurrently without any locking.
// There is NO log entry struct in gommon/log, which is used in many logging packages, the reason is
// if extra field is added to the interface, compiler would throw error on stale handler implementations.
type Handler interface {
	// HandleLog requires level, now, msg, all the others are optional
	// source is file:line, i.e. main.go:18 TODO: pass frame instead of string so handler can use trace for error handling?
	// context are fields attached to the logger instance
	// fields are ad-hoc fields from logger method like DebugF(msg, fields)
	HandleLog(level Level, now time.Time, msg string, source string, context Fields, fields Fields)
	// Flush writes the buffer to underlying storage
	Flush()
}

// NOTE: since the interface of handler become so big, there is no way to use a handler func
// TODO: maybe we can trim down the interface by adding more if else inside handler and use the max parameters
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

var defaultHandler = NewIOHandler(os.Stderr)

// DefaultHandler returns the singleton defaultHandler instance, which logs to stderr in text format
func DefaultHandler() Handler {
	return defaultHandler
}

// MultiHandler creates a handler that duplicates the log to all the provided handlers, it runs in
// serial and don't handle panic
func MultiHandler(handlers ...Handler) Handler {
	return &multiHandler{handlers: handlers}
}

// IOHandler writes log to io.Writer, default handler is a IOHandler using os.Stderr
// TODO: rename to text handler, this gonna break many applications ...
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

func (h *IOHandler) HandleLog(level Level, time time.Time, msg string, source string, context Fields, fields Fields) {
	b := make([]byte, 0, 50+len(msg))
	// level
	b = append(b, level.String()...)
	// time
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	// source, optional
	if source != "" {
		b = append(b, ' ')
		b = append(b, source...)
	}
	// message
	b = append(b, ' ')
	b = append(b, msg...)
	// context
	if len(context) > 0 {
		b = append(b, ' ')
		b = formatFields(b, context)
	}
	// fields
	if len(fields) > 0 {
		b = formatFields(b, fields)
	}
	b = append(b, '\n')
	h.w.Write(b)
}

// Flush implements Handler interface
func (h *IOHandler) Flush() {
	if s, ok := h.w.(Syncer); ok {
		s.Sync()
	}
}

// ----------------- start of multi handler ---------------------------

var _ Handler = (*multiHandler)(nil)

// https://github.com/dyweb/gommon/issues/87
type multiHandler struct {
	handlers []Handler
}

func (m *multiHandler) HandleLog(level Level, now time.Time, msg string, source string, context Fields, fields Fields) {
	for _, h := range m.handlers {
		h.HandleLog(level, now, msg, source, context, fields)
	}
}

func (m *multiHandler) Flush() {
	for _, h := range m.handlers {
		h.Flush()
	}
}

// ----------------- end of multi handler ---------------------------

// ----------------- start of text format util ---------------------------

// it has an extra tailing space, which can be updated in place to a \n
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
	b = b[:len(b)-1] // remove trailing space
	return b
}

// ----------------- end of text format util ---------------------------
