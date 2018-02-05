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
		dlog.Int("num", 2),
		dlog.Str("foo", "bar"),
	})
	log.EnableSource()
	log.Info("show me the line")
	log.Infof("show the line %d", 2)
	log.InfoF("show the line", dlog.Fields{
		dlog.Int("num", 2),
		dlog.Str("foo", "bar"),
	})
	log.DisableSource()
	log.Info("no more line number")

	// TODO: panic and fatal
}
