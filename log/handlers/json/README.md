# JSON log handler

For machine

## TODO

- [ ] string are NOT escaped, both message and fields
  - i.e. broken JSON like `"msg":"this "specific file"", "file":""there is space.yml""`

## Usage

````go
package main

import (
	"os"
	"time"
	
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/json"
)

var log = dlog.NewApplicationLogger()

func main() {
	dlog.SetHandlerRecursive(log, json.New(os.Stderr))
	log.Info("hi")
	log.Infof("open file %s", "foo.yml")
	log.InfoF("open", dlog.Fields{
		dlog.Str("file", "foo.yml"),
		dlog.Int("mode", 0666),
	})
	log.Warn("I am yellow")
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Info("recovered", r)
			}
		}()
		log.Panic("Panic reason")
	}()
	dlog.SetLevelRecursive(log, dlog.DebugLevel)
	log.Debug("I will sleep for a while")
	time.Sleep(1 * time.Second)
	log.Fatal("I am red")
}
````

````text
{"l":"info","t":1518210625,"m":"hi"}
{"l":"info","t":1518210625,"m":"open file foo.yml"}
{"l":"info","t":1518210625,"m":"open","file":"foo.yml","mode":438}
{"l":"warn","t":1518210625,"m":"I am yellow"}
{"l":"panic","t":1518210625,"m":"Panic reason"}
{"l":"info","t":1518210625,"m":"recoveredPanic reason"}
{"l":"debug","t":1518210625,"m":"I will sleep for a while"}
{"l":"fatal","t":1518210626,"m":"I am red"}
````