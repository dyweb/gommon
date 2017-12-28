# gokit-log

https://github.com/go-kit/kit/tree/master/log

- complexity: low

[Slide: The Hunt for a Logger Interface](http://go-talks.appspot.com/github.com/ChrisHines/talks/structured-logging/structured-logging.slide)

- mainly about design interface, logging is just the context
- glog is ad-hoc, non structual
- log15 allow context `log.Info("I am msg", "t1", "v1", "t2", "v2")`, `clog = log.New("common", "field")`
- mainly focus on small interface (just one method, seems to extreme to me)


- it supports redirect stdlib logger
  - stdlib.go, `NewStdlibAdapter` which implements `Write([] byte)` and use regexp to make it structual ...
  - you can also do the reverse ....
- use `github.com/go-stack/stack` to obtain caller (file, line number etc.)
- level is implemented in `level` subpackage, it wraps a Logger

````go
// log.go
type Logger interface {
	Log(keyvals ...interface{}) error
}

// level/level.go
type logger struct {
	next           log.Logger
	allowed        level
	squelchNoLevel bool
	errNotAllowed  error
	errNoLevel     error
}
````

````go
w := log.NewSyncWriter(os.Stderr)
logger := log.NewLogfmtLogger(w)
logger.Log("question", "what is the meaning of life?", "answer", 42)

// Output:
// question="what is the meaning of life?" answer=42
````