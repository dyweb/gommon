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
	// HandleLog accepts level, log time, formatted log message
	HandleLog(level Level, now time.Time, msg string)
	// HandleLogWithSource accepts formatted source line of log i.e., http.go:13
	// TODO: pass frame instead of string so handler can use trace for error handling?
	HandleLogWithSource(level Level, now time.Time, msg string, source string)
	// HandleLogWithFields accepts fields with type hint,
	// implementation should inspect the type field instead of using reflection
	HandleLogWithFields(level Level, now time.Time, msg string, fields Fields)
	// HandleLogWithSourceFields accepts both source and fields
	HandleLogWithSourceFields(level Level, now time.Time, msg string, source string, fields Fields)
	// HandleLogWithContextFields get context from logger, which is also fields
	HandleLogWithContextFields(level Level, now time.Time, msg string, context Fields, fields Fields)
	// HandleLogWithSourceContextFields contains everything
	HandleLogWithSourceContextFields(level Level, now time.Time, msg string, source string, context Fields, fields Fields)
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
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = formatHead(b, level, time, msg)
	b = append(b, '\n')
	h.w.Write(b)
}

// HandleLogWithSource implements Handler interface
func (h *IOHandler) HandleLogWithSource(level Level, time time.Time, msg string, source string) {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg)+len(source))
	b = formatHeadWithSource(b, level, time, msg, source)
	b = append(b, '\n')
	h.w.Write(b)
}

// HandleLogWithFields implements Handler interface
func (h *IOHandler) HandleLogWithFields(level Level, time time.Time, msg string, fields Fields) {
	// we use raw slice instead of bytes buffer because we need to use strconv.Append*, which requires raw slice
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = formatHead(b, level, time, msg)
	b = append(b, ' ')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

// HandleLogWithSourceFields implements Handler interface
func (h *IOHandler) HandleLogWithSourceFields(level Level, time time.Time, msg string, source string, fields Fields) {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg)+len(source))
	b = formatHeadWithSource(b, level, time, msg, source)
	b = append(b, ' ')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

func (h *IOHandler) HandleLogWithContextFields(level Level, time time.Time, msg string, context Fields, fields Fields) {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = formatHead(b, level, time, msg)
	b = append(b, ' ')
	b = formatFields(b, context)
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

func (h *IOHandler) HandleLogWithSourceContextFields(level Level, time time.Time, msg string, source string, context Fields, fields Fields) {
	b := make([]byte, 0, 5+4+len(defaultTimeStampFormat)+len(msg))
	b = formatHeadWithSource(b, level, time, msg, source)
	b = append(b, ' ')
	b = formatFields(b, context)
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

// ----------------- start of multi handler ---------------------------

var _ Handler = (*multiHandler)(nil)

// https://github.com/dyweb/gommon/issues/87
type multiHandler struct {
	handlers []Handler
}

func (m *multiHandler) HandleLog(level Level, now time.Time, msg string) {
	for _, h := range m.handlers {
		h.HandleLog(level, now, msg)
	}
}

func (m *multiHandler) HandleLogWithSource(level Level, now time.Time, msg string, source string) {
	for _, h := range m.handlers {
		h.HandleLogWithSource(level, now, msg, source)
	}
}

func (m *multiHandler) HandleLogWithFields(level Level, now time.Time, msg string, fields Fields) {
	for _, h := range m.handlers {
		h.HandleLogWithFields(level, now, msg, fields)
	}
}

func (m *multiHandler) HandleLogWithSourceFields(level Level, now time.Time, msg string, source string, fields Fields) {
	for _, h := range m.handlers {
		h.HandleLogWithSourceFields(level, now, msg, source, fields)
	}
}

func (m *multiHandler) HandleLogWithContextFields(level Level, now time.Time, msg string, context Fields, fields Fields) {
	for _, h := range m.handlers {
		h.HandleLogWithContextFields(level, now, msg, context, fields)
	}
}

func (m *multiHandler) HandleLogWithSourceContextFields(level Level, now time.Time, msg string, source string, context Fields, fields Fields) {
	for _, h := range m.handlers {
		h.HandleLogWithSourceContextFields(level, now, msg, source, context, fields)
	}
}

func (m *multiHandler) Flush() {
	for _, h := range m.handlers {
		h.Flush()
	}
}

// ----------------- end of multi handler ---------------------------

// ----------------- start of text format util ---------------------------

// no need to use fmt.Printf since we don't need any format
func formatHead(b []byte, level Level, time time.Time, msg string) []byte {
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, msg...)
	return b
}

// we have a new function because source sits between time and msg in output, instead of after msg
// i.e. info 2018-02-04T21:03:20-08:00 main.go:18 show me the line
func formatHeadWithSource(b []byte, level Level, time time.Time, msg string, source string) []byte {
	b = append(b, level.String()...)
	b = append(b, ' ')
	b = time.AppendFormat(b, defaultTimeStampFormat)
	b = append(b, ' ')
	b = append(b, source...)
	b = append(b, ' ')
	b = append(b, msg...)
	return b
}

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
	return b
}

// ----------------- end of text format util ---------------------------
