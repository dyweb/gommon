package log

import (
	"io"
	"os"
	"strconv"
	"strings"
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
	// source is Caller which contains full file line TODO: pass frame instead of string so handler can use trace for error handling?
	// context are fields attached to the logger instance
	// fields are ad-hoc fields from logger method like DebugF(msg, fields)
	HandleLog(level Level, now time.Time, msg string, source Caller, context Fields, fields Fields)
	// Flush writes the buffer to underlying storage
	Flush()
}

// HandlerFunc is an adapter to allow use of ordinary functions as log entry handlers
type HandlerFunc func(level Level, now time.Time, msg string, source string, context Fields, fields Fields)

// TODO: why the receiver is value instead of pointer https://github.com/dyweb/gommon/issues/30 and what's the overhead
func (f HandlerFunc) HandleLog(level Level, now time.Time, msg string, source string, context Fields, fields Fields) {
	f(level, now, msg, source, context, fields)
}

var _ Syncer = (*os.File)(nil)

// Syncer is implemented by os.File, handler implementation should check this interface and call Sync
// if they support using file as sink
// TODO: about sync
// - in go, both os.Stderr and os.Stdout are not (line) buffered
// - what would happen if os.Stderr.Close()
// - os.Stderr.Sync() will there be any different if stderr/stdout is redirected to a file
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

// TextHandler writes log to io.Writer, default handler is a TextHandler using os.Stderr
type TextHandler struct {
	w io.Writer
}

// NewIOHandler returns the default text handler, the name is for backward compatibility
func NewIOHandler(w io.Writer) Handler {
	return &TextHandler{w: w}
}

// NewTextHandler formats log in human readable format without any escape, thus it is NOT machine readable
func NewTextHandler(w io.Writer) Handler {
	return &TextHandler{w: w}
}

func (h *TextHandler) HandleLog(level Level, time time.Time, msg string, source Caller, context Fields, fields Fields) {
	b := make([]byte, 0, 50+len(msg)+len(source.File)+30*len(context)+30*len(fields))
	// level
	b = append(b, level.AlignedUpperString()...)
	// time
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	// source, optional
	if source.Line != 0 {
		b = append(b, ' ')
		last := strings.LastIndex(source.File, "/")
		b = append(b, source.File[last+1:]...)
		b = append(b, ':')
		b = strconv.AppendInt(b, int64(source.Line), 10)
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
func (h *TextHandler) Flush() {
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

func (m *multiHandler) HandleLog(level Level, now time.Time, msg string, source Caller, context Fields, fields Fields) {
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
