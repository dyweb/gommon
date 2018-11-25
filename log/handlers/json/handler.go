/*
Package json writes log in JSON format, it concatenates string directly and does not use encoding/json.
TODO: support escape string, and compare with standard json encoding
*/
package json

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

func (h *Handler) HandleLog(level log.Level, time time.Time, msg string, source string, context log.Fields, fields log.Fields) {
	// FIXME: the hard coded 50 is not correct, it depends on source and fields etc.
	// TODO: I think this make can only be optimized by using a buffer pool and
	b := make([]byte, 0, 50+len(msg))
	// level
	b = append(b, `{"l":"`...)
	b = append(b, level.String()...)
	// time
	b = append(b, `","t":`...)
	b = strconv.AppendInt(b, time.Unix(), 10)
	// message
	b = append(b, `,"m":"`...)
	b = append(b, msg...)
	b = append(b, '"')
	// source
	if source != "" {
		b = append(b, `,"s":"`...)
		b = append(b, source...)
		b = append(b, '"')
	}
	// context
	if len(context) > 0 {
		b = append(b, `,`...)
		b = formatFields(b, context)
	}
	// fields
	if len(fields) > 0 {
		b = append(b, `,`...)
		b = formatFields(b, fields)
	}
	b = append(b, "}\n"...)
	h.w.Write(b)
}

func (h *Handler) Flush() {
	if s, ok := h.w.(log.Syncer); ok {
		s.Sync()
	}
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
