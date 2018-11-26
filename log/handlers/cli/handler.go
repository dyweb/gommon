/*
Package cli writes is same as builtin IOHandler except color and delta time.
It is used by go.ice as default handler
TODO: color can't be disabled and we don't detect tty like logrus
*/
package cli

import (
	"io"
	"strconv"
	"time"

	"github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/util/color"
)

const (
	defaultTimeStampFormat = time.RFC3339
)

var _ log.Handler = (*Handler)(nil)

type Handler struct {
	w     io.Writer
	start time.Time
	delta bool
}

func New(w io.Writer, delta bool) *Handler {
	return &Handler{
		w:     w,
		start: time.Now(),
		delta: delta,
	}
}

func (h *Handler) HandleLog(level log.Level, now time.Time, msg string, source string, context log.Fields, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg)+len(source)+30*len(context)+30*len(fields))
	// level
	b = append(b, level.ColoredAlignedUpperString()...)
	// time
	b = append(b, ' ')
	if h.delta {
		b = append(b, formatNum(uint(now.Sub(h.start)/time.Second), 4)...)
	} else {
		b = now.AppendFormat(b, defaultTimeStampFormat)
	}
	// source
	if source != "" {
		b = append(b, ' ')
		b = append(b, color.CyanStart...)
		b = append(b, source...)
		b = append(b, color.End...)
	}
	// message
	b = append(b, ' ')
	b = append(b, msg...)
	// context
	if len(context) > 0 {
		b = append(b, ' ')
		b = formatFields(b, context)
	}
	// field
	if len(fields) > 0 {
		b = append(b, ' ')
		b = formatFields(b, fields)
	}
	b = append(b, '\n')
	h.w.Write(b)
}

func (h *Handler) Flush() {
	if s, ok := h.w.(log.Syncer); ok {
		s.Sync()
	}
}

// NOTE: most of format functions are same as (copied from) default handler

// based on fmt.fmt_integer, only support base 10 and prefix 0
// https://golang.org/src/fmt/format.go#L194
func formatNum(u uint, digits int) []byte {
	i := digits
	b := make([]byte, digits)
	for u > 0 && i > 0 {
		i--
		next := u / 10
		b[i] = byte('0' + u - next*10)
		u = next
	}
	for i > 0 {
		i--
		b[i] = '0'
	}
	return b
}

// it has an extra tailing space, which can be updated inplace to a \n
func formatFields(b []byte, fields log.Fields) []byte {
	for _, f := range fields {
		b = append(b, color.CyanStart...)
		b = append(b, f.Key...)
		b = append(b, color.End...)
		b = append(b, '=')
		switch f.Type {
		case log.IntType:
			b = strconv.AppendInt(b, f.Int, 10)
		case log.StringType:
			b = append(b, f.Str...)
		}
		b = append(b, ' ')
	}
	b = b[:len(b)-1] // remove trailing space
	return b
}
