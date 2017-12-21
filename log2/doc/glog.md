# glog

https://github.com/golang/glog

- complexity: medium
- supports rotate file and log to different file?
- log methods like `Info`, `Infof` are implemented on `Verbose` which is a boolean indicate if required level should be logged
- includes buffer etc.
- lock when log
- it **collects counter** (needed for fancy features like when logging hits line file:N emit a stack trace)

````go
if glog.V(2) {
    glog.Info("Starting transaction...")
}

glog.V(2).Infoln("Processed", nItems, "elements")
````

````go
func (l *loggingT) output(s severity, buf *buffer, file string, line int, alsoToStderr bool) {
	l.mu.Lock()
	// ...
    l.mu.Unlock()
    if stats := severityStats[s]; stats != nil {
        atomic.AddInt64(&stats.lines, 1)
        atomic.AddInt64(&stats.bytes, int64(len(data)))
    }
}
````

flags, including interesting features

````go
func init() {
	flag.BoolVar(&logging.toStderr, "logtostderr", false, "log to standard error instead of files")
    flag.BoolVar(&logging.alsoToStderr, "alsologtostderr", false, "log to standard error as well as files")
    flag.Var(&logging.verbosity, "v", "log level for V logs")
    flag.Var(&logging.stderrThreshold, "stderrthreshold", "logs at or above this threshold go to stderr")
    flag.Var(&logging.vmodule, "vmodule", "comma-separated list of pattern=N settings for file-filtered logging")
    flag.Var(&logging.traceLocation, "log_backtrace_at", "when logging hits line file:N, emit a stack trace")
}
````

The cpp version actually supports more features 

https://github.com/google/glog

`LOG_IF(INFO, num_cookies > 10)`, `LOG_EVERY_N(INFO, 10)`, , `LOG_IF_EVERY_N`, `LOG_FIRST_N(INFO, 20)`

- macro is used to expend those to real function calls, it is using automake `logging.h.in` so