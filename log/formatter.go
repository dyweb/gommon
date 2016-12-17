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
	// TODO: print fields
	fmt.Fprintf(b, "[%s]%s\n", entry.Level.String(), entry.Message)
	return b.Bytes(), nil
}
