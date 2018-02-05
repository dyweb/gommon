package main

import (
	dlog "github.com/dyweb/gommon/log"
)

var log = dlog.NewApplicationLogger()

// simply log to stderr
func main() {
	log.Info("this is love!")
	log.Infof("this is love %d", 2)
	log.InfoF("this love", dlog.Fields{
		dlog.Field{Key: "num", Value: 2},
	})
}
