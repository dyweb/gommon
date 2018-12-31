// +build ignore

package log

import (
	"bytes"
	"fmt"
	"time"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
	SetColor(bool)
	SetElapsedTime(bool)
	SetTimeFormat(string)
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

const (
	defaultTimeStampFormat = time.RFC3339
)

var _ Formatter = (*TextFormatter)(nil)

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
		TimeStampFormat:   defaultTimeStampFormat,
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
			// https://github.com/sirupsen/logrus/pull/465
			fmt.Fprintf(b, "[%04d]", int(entry.Time.Sub(baseTimeStamp)/time.Second))
		} else {
			// NOTE: config.Validate would check if timestamp format set by user is valid, and time.Format does not return error
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

func (f *TextFormatter) SetColor(b bool) {
	f.EnableColor = b
}

func (f *TextFormatter) SetElapsedTime(b bool) {
	f.EnableElapsedTime = b
}

func (f *TextFormatter) SetTimeFormat(tf string) {
	if tf == "" {
		f.TimeStampFormat = defaultTimeStampFormat
	} else {
		f.TimeStampFormat = tf
	}
}
