package main

import (
	"github.com/dyweb/gommon/log2"
	"github.com/dyweb/gommon/log2/handlers"
)

func main() {
	logger := log2.Logger{}
	logger.SetLevel(log2.DebugLevel)
	logger.SetHandler(handlers.NewStdout())
	// FIXME: extra brackets [] and no new line
	// debug [This is a debug message 1]info [This is a info message 2]
	logger.Debug("This is a debug message", 1)
	logger.Info("This is a info message", 2)
}
