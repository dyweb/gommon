/*
Package json writes log in JSON format, it concatenates string directly and does not use encoding/json.
TODO: support escape
*/
package json // import "github.com/dyweb/gommon/log/handlers/json"

import (
	"io"
	"strconv"
	"time"

	"github.com/dyweb/gommon/log"
)

var _ log.Handler = (*Handler)(nil)

type Handler struct {
	w io.Writer
}

func New(w io.Writer) *Handler {
	return &Handler{
		w: w,
	}
}

func (h *Handler) HandleLog(level log.Level, time time.Time, msg string) {
	// FIXME: the hard coded 50 is not correct, it depends on source and fields etc.
	// TODO: I think this make can only be optimized by using a buffer pool and
	// we should have buffer pool for different
	b := make([]byte, 0, 50+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithSource(level log.Level, time time.Time, msg string, source string) {
	b := make([]byte, 0, 50+len(source)+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, `,"s":"`...)
	b = append(b, source...)
	b = append(b, "\"}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithFields(level log.Level, time time.Time, msg string, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, ',')
	b = formatFields(b, fields)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithSourceFields(level log.Level, time time.Time, msg string, source string, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, `,"s":"`...)
	b = append(b, source...)
	b = append(b, `",`...)
	b = formatFields(b, fields)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithContextFields(level log.Level, time time.Time, msg string, context log.Fields, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, ',')
	b = formatFields(b, context)
	b = append(b, ',')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

func (h *Handler) HandleLogWithSourceContextFields(level log.Level, time time.Time, msg string, source string, context log.Fields, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg))
	b = formatHead(b, level, time.Unix(), msg)
	b = append(b, `,"s":"`...)
	b = append(b, source...)
	b = append(b, `",`...)
	b = formatFields(b, context)
	b = append(b, ',')
	b = formatFields(b, fields)
	b[len(b)-1] = '\n'
	h.w.Write(b)
}

func (h *Handler) Flush() {
	if s, ok := h.w.(log.Syncer); ok {
		s.Sync()
	}
}

func formatHead(b []byte, level log.Level, time int64, msg string) []byte {
	b = append(b, `{"l":"`...)
	b = append(b, level.String()...)
	b = append(b, `","t":`...)
	b = strconv.AppendInt(b, time, 10)
	b = append(b, `,"m":"`...)
	b = append(b, msg...)
	b = append(b, '"')
	return b
}

func formatFields(b []byte, fields log.Fields) []byte {
	for _, f := range fields {
		b = append(b, '"')
		b = append(b, f.Key...)
		b = append(b, "\":"...)
		switch f.Type {
		case log.IntType:
			b = strconv.AppendInt(b, f.Int, 10)
		case log.StringType:
			b = append(b, '"')
			b = append(b, f.Str...)
			b = append(b, '"')
		}
		b = append(b, ',')
	}
	b = b[:len(b)-1] // remove trailing comma
	return b
}
