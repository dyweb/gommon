// Package cli generates human readable text with color and display time in delta.
// Color and delta can be disabled and output will be same as default handler
// It does NOT escape quote so it is not machine readable
//
//	func main() {
//		var log, logReg = dlog.NewApplicationLoggerAndRegistry("example")
// 		dlog.SetHandler(logReg, cli.New(os.Stderr, true)) // with color, delta time
//	}
//
package cli

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/util/color"
)

const (
	DefaultTimeStampFormat = time.RFC3339
)

var _ log.Handler = (*handler)(nil)

// handler is not exported since all the fields are not exported
type handler struct {
	w     io.Writer
	start time.Time
	delta bool
	color bool
}

// New returns a handler with color on level and field key
func New(w io.Writer, delta bool) log.Handler {
	return &handler{
		w:     w,
		start: time.Now(),
		delta: delta,
		color: true,
	}
}

// NewNoColor returns a handler without color and use full timestamp
func NewNoColor(w io.Writer) log.Handler {
	return &handler{
		w:     w,
		start: time.Now(),
		delta: false,
		color: false,
	}
}

func (h *handler) HandleLog(level log.Level, now time.Time, msg string, source log.Caller, context log.Fields, fields log.Fields) {
	b := make([]byte, 0, 50+len(msg)+len(source.File)+30*len(context)+30*len(fields))
	// level
	if h.color {
		b = append(b, level.ColoredAlignedUpperString()...)
	} else {
		b = append(b, level.AlignedUpperString()...)
	}
	// time
	b = append(b, ' ')
	if h.delta {
		b = append(b, formatNum(uint(now.Sub(h.start)/time.Second), 4)...)
	} else {
		b = now.AppendFormat(b, DefaultTimeStampFormat)
	}
	// source
	if source.Line != 0 {
		b = append(b, ' ')
		if h.color {
			b = append(b, color.CyanStart...)
		}
		last := strings.LastIndex(source.File, "/")
		b = append(b, source.File[last+1:]...)
		b = append(b, ':')
		b = strconv.AppendInt(b, int64(source.Line), 10)
		if h.color {
			b = append(b, color.End...)
		}
	}
	// message
	b = append(b, ' ')
	b = append(b, msg...)
	// context
	if len(context) > 0 {
		b = append(b, ' ')
		b = formatFields(b, h.color, context)
	}
	// field
	if len(fields) > 0 {
		b = append(b, ' ')
		b = formatFields(b, h.color, fields)
	}
	b = append(b, '\n')
	h.w.Write(b)
}

func (h *handler) Flush() {
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

func formatFields(b []byte, useColor bool, fields log.Fields) []byte {
	for _, f := range fields {
		if useColor {
			b = append(b, color.CyanStart...)
		}
		b = append(b, f.Key...)
		if useColor {
			b = append(b, color.End...)
		}
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
