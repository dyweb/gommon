/*
Package json writes log in JSON format, it concatenates string directly and does not use encoding/json.
TODO: support escape
*/
package json // import "github.com/dyweb/gommon/log/handlers/json"

import (
	"github.com/dyweb/gommon/log"
	"io"
	"strconv"
	"time"
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
	b := formatHead(level, time, msg)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithSource(level log.Level, time time.Time, msg string, source string) {
	b := formatHead(level, time, msg)
	b = append(b, `,"s":"`...)
	b = append(b, source...)
	b = append(b, "\"}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithFields(level log.Level, time time.Time, msg string, fields log.Fields) {
	b := formatHead(level, time, msg)
	b = append(b, ',')
	b = formatFields(b, fields)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) HandleLogWithSourceFields(level log.Level, time time.Time, msg string, source string, fields log.Fields) {
	b := formatHead(level, time, msg)
	b = append(b, `,"s":"`...)
	b = append(b, source...)
	b = append(b, `",`...)
	b = formatFields(b, fields)
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) Flush() {
	if s, ok := h.w.(log.Syncer); ok {
		s.Sync()
	}
}

func formatHead(level log.Level, time time.Time, msg string) []byte {
	b := make([]byte, 0, 5+4+10+len(msg))
	b = append(b, `{"l":"`...)
	b = append(b, level.String()...)
	b = append(b, `","t":`...)
	b = strconv.AppendInt(b, time.Unix(), 10)
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
