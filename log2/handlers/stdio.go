package handlers

import (
	"fmt"
	"io"
	"os"

	"github.com/dyweb/gommon/log2"
)

var _ log2.Handler = (*StdIO)(nil)

type StdIO struct {
	// TODO: does stdout/stderr has Close method?
	w io.WriteCloser
}

func NewStdout() *StdIO {
	return &StdIO{
		w: os.Stdout,
	}
}

func (s *StdIO) HandleLog(level log2.Level, msg string) {
	// TODO: there should be more efficient way of coverting level to string, not using fmt.Stringer
	// TODO: might use a buffer, since we no longer need format of fmt
	fmt.Fprintf(s.w, "%s %s\n", level, msg)
}

func (s *StdIO) Flush() {
	s.w.Close()
}
