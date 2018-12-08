package main

import (
	"os"
	"time"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/cli"
	"github.com/dyweb/gommon/log/handlers/json"
)

var log, logReg = dlog.NewApplicationLoggerAndRegistry("simple")

// simply log to stderr
func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "json" {
			dlog.SetHandler(logReg, json.New(os.Stderr))
		}
		if os.Args[1] == "cli" {
			dlog.SetHandler(logReg, cli.New(os.Stderr, false))
		}
		if os.Args[1] == "cli-d" {
			dlog.SetHandler(logReg, cli.New(os.Stderr, true))
		}
	}
	dlog.SetLevel(logReg, dlog.DebugLevel)
	log.Debug("show me the meaning of being lonely")
	log.Info("this is love!")
	log.Print("print is info level")
	log.Warnf("this is love %d", 2)
	log.InfoF("this love", dlog.Int("num", 2), dlog.Str("foo", "bar"))
	dlog.EnableSource(logReg)
	// TODO: show skip caller using a util func
	log.Info("show me the line")
	log.Infof("show the line %d", 2)
	log.InfoF("show the line", dlog.Int("num", 2), dlog.Str("foo", "bar"))
	log.DisableSource()
	log.WarnF("I will sleep", dlog.Int("duration", 1))
	time.Sleep(10 * time.Millisecond)
	log.Info("no more line number")

	log.AddField(dlog.Str("f1", "v1"))
	log.Info("should have some extra context")
	// TODO: panic and fatal
}
