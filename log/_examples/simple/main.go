package main

import (
	"github.com/dyweb/gommon/log2"
	"github.com/dyweb/gommon/log2/handlers"
)

// simply log to stdout
// TODO: should run examples when test
func main() {
	logger := log2.Logger{}
	logger.SetLevel(log2.DebugLevel)
	logger.SetHandler(handlers.NewStdout())
	logger.Debug("This is a debug message ", 1, " yeah")
	logger.Info("This is a info message ", 2, " no")
}
