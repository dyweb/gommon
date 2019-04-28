# JSON log handler

For machine

## TODO

- [ ] allow config built in filed names

## Usage

See [example/main.go](example/main.go)

TODO: gommon should allow sync file into markdown

````go
package main

import (
	"os"
	"time"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/json"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()

func main() {
	dlog.SetHandler(json.New(os.Stderr))
	log.Info("hi")
	log.Infof("open file %s", "foo.yml")
	log.InfoF("open",
		dlog.Str("file", "foo.yml"),
		dlog.Int("mode", 0666),
	)
	log.Warn("I am yellow")
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Info("recovered", r)
			}
		}()
		log.Panic("I just want to panic")
	}()
	dlog.SetLevel(dlog.DebugLevel)
	log.Debug("I will sleep for a while")
	time.Sleep(500 * time.Millisecond)
	log.Fatal("I am red")
}
````

````text
{"l":"info","t":1546234314,"m":"hi"}
{"l":"info","t":1546234314,"m":"open file foo.yml"}
{"l":"info","t":1546234314,"m":"open","file":"foo.yml","mode":438}
{"l":"warn","t":1546234314,"m":"I am yellow"}
{"l":"panic","t":1546234314,"m":"I just want to panic","s":"main.go:28"}
{"l":"info","t":1546234314,"m":"recoveredI just want to panic"}
{"l":"debug","t":1546234314,"m":"I will sleep for a while"}
{"l":"fatal","t":1546234314,"m":"I am red","s":"main.go:33"}
````