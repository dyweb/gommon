package main

import (
	"github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers"
)

// simply log to stdout
// TODO: should run examples when test
func main() {
	logger := log.Logger{}
	logger.SetLevel(log.DebugLevel)
	logger.SetHandler(handlers.NewStdout())
	logger.Debug("This is a debug message ", 1, " yeah")
	logger.Info("This is a info message ", 2, " no")
}
