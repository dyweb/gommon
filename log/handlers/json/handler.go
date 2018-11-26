/*
Package json writes log in JSON format, it concatenates string directly and does not use encoding/json.
TODO: compare with standard json encoding
*/
package json

import (
	"io"
	"strconv"
	"time"
	"unicode/utf8"

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
	b = append(b, `,"m":`...)
	b = encodeString(b, msg)
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
		// TODO: should we also escape key? ...
		b = append(b, '"')
		b = append(b, f.Key...)
		b = append(b, "\":"...)
		switch f.Type {
		case log.IntType:
			b = strconv.AppendInt(b, f.Int, 10)
		case log.StringType:
			b = encodeString(b, f.Str)
		}
		b = append(b, ',')
	}
	b = b[:len(b)-1] // remove trailing comma
	return b
}

// encodeString escape character like " \n, it does not handle jsonp or html like standard library does
// it is based on encoding/json/encode.go func (e *encodeState) string(s string, escapeHTML bool) w/ some comment
func encodeString(buf []byte, s string) []byte {
	buf = append(buf, '"')
	start := 0
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf { // characters below RuneSelf are represented as themselves in a single byte.
			if safeSet[b] {
				i++
				// we don't call append right away for this byte
				// because if the entire string is safe, we only need to call append once
				continue
			}
			// append previous processed bytes
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			// some special ascii characters need escape
			switch b {
			case '\\', '"':
				buf = append(buf, '\\', b)
			case '\n':
				buf = append(buf, '\\', 'n')
			case '\r':
				buf = append(buf, '\\', 'r')
			case '\t':
				buf = append(buf, '\\', 't')
			default:
				// TODO: I don't get what this section does ...
				buf = append(buf, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		// it's utf8 rune
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			// when error, first append previous processed bytes
			if start < i {
				buf = append(buf, s[start:i]...)
			}
			buf = append(buf, `\ufffd`...)
			i += size
			start = i
			continue
		}
		i += size // only move the cursor, append happens when need to escape or there is error or at last
	}
	if start < len(s) {
		// NOTE: it's s[start:] not buf[star:] ....
		buf = append(buf, s[start:]...)
	}
	buf = append(buf, '"')
	return buf
}
