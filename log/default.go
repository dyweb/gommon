package log

import (
	"fmt"
	"os"
)

type defaultHandler struct {
}

var DefaultHandler Handler = &defaultHandler{}

func (h *defaultHandler) HandleLog(level Level, msg string) {
	// TODO: might use a buffer, since we are just concat string, no special format is needed
	// TODO: time, fields etc.
	fmt.Fprintf(os.Stdout, "%s %s\n", level, msg)
}

func (h *defaultHandler) Flush() {
	// nop
}
