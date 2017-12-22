package handlers

import "io"

type Std struct {
	// TODO: does stdout/stderr has Close method?
	writer io.WriteCloser
}

func (s *Std) HandleLog() {
	s.writer.Write([]byte(""))
}

func (s *Std) Flush() {
	s.writer.Close()
}
