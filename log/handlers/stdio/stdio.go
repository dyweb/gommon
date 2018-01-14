// Package stdio can be used with stderr/stdout/file and all the writer that implemented WriteSyncer interface
// its code is almost identical with default handler, but we duplicate the code to avoid cycle import
package stdio

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dyweb/gommon/log"
)

// interface check
var _ log.Handler = (*Handler)(nil)

const (
	defaultTimeStampFormat = time.RFC3339
)

// WriteSyncer is implemented by os.File thus os.Stdout and os.Stdout
// it is based on uber/zap/zapcore/write_syncer.go
type WriteSyncer interface {
	io.Writer
	Sync() error
}

type Handler struct {
	// TODO: Close or not? it seems for os.Stdout and os.Stderr, it should not be closed
	w        WriteSyncer
	tsFormat string
}

func NewStdout() *Handler {
	h := newstd()
	h.w = os.Stdout
	return h
}

func NewStderr() *Handler {
	h := newstd()
	h.w = os.Stderr
	return h
}

// TODO: append or write etc.
//func NewFile() (*Handler, error) {
//
//}

func newstd() *Handler {
	return &Handler{
		tsFormat: defaultTimeStampFormat,
	}
}

// TODO: sync w/ default handler on fields etc.
func (s *Handler) HandleLog(level log.Level, msg string) {
	fmt.Fprintf(s.w, "%s %s\n", level.String(), msg)
	fmt.Fprintf(os.Stderr, "%s %s %s\n", level.String(), time.Now().Format(s.tsFormat), msg)
}

func (s *Handler) Flush() {
	s.w.Sync()
}
