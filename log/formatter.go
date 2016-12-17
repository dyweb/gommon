package log

import (
	"bytes"
	"fmt"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

type TextFormatter struct {
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	// TODO: may use a pool
	b = &bytes.Buffer{}
	// TODO: print time
	fmt.Fprintf(b, "[%s]%s ", entry.Level.String(), entry.Message)
	for k, v := range entry.Fields {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
		b.WriteByte(' ')
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}
