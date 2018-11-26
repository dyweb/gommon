package main

import (
	"io/ioutil"

	dlog "github.com/dyweb/gommon/log"
)

//
//type myLogger struct {
//	h handler
//}
//
////go:noinline
//func (l *myLogger) info(s string) {
//	l.h.log(s)
//}
//
//type handler interface {
//	log(s string)
//}
//
//type printer struct {
//}
//
////go:noinline
//func (p *printer) log(s string) {
//	// do nothing
//	os.Stdout.Write([]byte(s))
//}

// go build -gcflags "-m -m" .
func main() {
	//mLogger := myLogger{
	//	h: &printer{},
	//}
	//mLogger.info("a")

	logger := dlog.NewTestLogger(dlog.InfoLevel)
	logger.SetHandler(dlog.NewIOHandler(ioutil.Discard))

	logger.InfoF("a")                   // no slice of fields, no heap alloc
	logger.InfoF("a", dlog.Int("a", 1)) // escaped
	logger.NoopF("a", dlog.Int("a", 1)) // NoopF don't call any interface using the fields given, so no heap
}
