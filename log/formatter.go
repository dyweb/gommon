package log

import (
	"bytes"
	"fmt"
	"time"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

var (
	baseTimeStamp = time.Now()
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

type TextFormatter struct {
	EnableColor       bool
	EnableTimeStamp   bool
	EnableElapsedTime bool
	TimeStampFormat   string
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		EnableColor:       false,
		EnableTimeStamp:   true,
		EnableElapsedTime: true,
		TimeStampFormat:   time.RFC3339,
	}
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	// TODO: may use a pool
	b = &bytes.Buffer{}
	var levelColor = nocolor
	if f.EnableColor {
		switch entry.Level {
		case InfoLevel:
			levelColor = blue
		case WarnLevel:
			levelColor = yellow
		case ErrorLevel, FatalLevel, PanicLevel:
			levelColor = red
		}
	}
	if levelColor != nocolor {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", levelColor, entry.Level.ShortUpperString())
	} else {
		b.WriteString(entry.Level.ShortUpperString())
	}
	if f.EnableTimeStamp {
		if f.EnableElapsedTime {
			// NOTE: the elapsedTime copied from logrus is wrong
			fmt.Fprintf(b, "[%04d]", int(entry.Time.Sub(baseTimeStamp)/time.Second))
		} else {
			// TODO: what if the user set TimeStampFormat to an invalid format
			fmt.Fprintf(b, "[%s]", entry.Time.Format(f.TimeStampFormat))
		}
	}
	b.WriteByte(' ')
	b.WriteString(entry.Message)
	b.WriteByte(' ')
	for k, v := range entry.Fields {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
		b.WriteByte(' ')
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

// NOTE: the elapsedTime copied from logrus is wrong
func elapsedTime() int {
	return int(time.Since(baseTimeStamp) / time.Second)
}
