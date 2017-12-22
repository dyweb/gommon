package handlers

import (
	"io"
	"github.com/dyweb/gommon/log2"
	"fmt"
	"os"
)

var _ log2.Handler = (*Std)(nil)

type Std struct {
	// TODO: does stdout/stderr has Close method?
	w io.WriteCloser
}

func NewStdout() *Std {
	return &Std{
		w: os.Stdout,
	}
}

func (s *Std) HandleLog(level log2.Level, msg string) {
	// TODO: there should be more efficient way of coverting level to string, not using fmt.Stringer
	// TODO: might use a buffer, since we no longer need format of fmt
	// FIXME: the message is debug [This is a debug message 1]info [This is a info message 2]
	fmt.Fprintf(s.w, "%s %s", level, msg)
}

func (s *Std) Flush() {
	s.w.Close()
}
