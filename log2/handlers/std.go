package handlers

import (
	"fmt"
	"io"
	"os"

	"github.com/dyweb/gommon/log2"
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
	fmt.Fprintf(s.w, "%s %s\n", level, msg)
}

func (s *Std) Flush() {
	s.w.Close()
}
