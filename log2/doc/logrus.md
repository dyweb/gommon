# Logrus

https://github.com/sirupsen/logrus

The first version of gommon/log is modeled after logrus, also found [a small bug #457 in its timing output when hook is involved](https://github.com/sirupsen/logrus/issues/457)

````text
alt_exit.go run handlers before call os.Exit
entry.go contain fields and message and all the logging method Debug, Info, Warn, Error etc.
formatter.go interface, implementation is in json_formatter and text_formatter
hooks.go interface, a hook can be attached to different level
json_formatter.go use encoding/json
logger.go io.Writer, formatter, level, mutex, pool for entry, all the methods a Entry has (not in gommon)
logrus.go logger interface, levels
terminal_*.go different OS (appengine)
text_formatter.go use a bytes buffer (Entry struct also has a pointer to buffer)
writer.go return *io.PipeWriter, so it can be used as normal writer
````

- its `Entry.log` is using non pointer receiver to avoid race condition
  - [at15/go-learning#3](https://github.com/at15/go-learning/issues/3)
  - [dyweb/Ayi#59](https://github.com/dyweb/Ayi/issues/59)
  
````go
// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	entry.Logger.mu.Lock()
	err := entry.Logger.Hooks.Fire(level, &entry)
	entry.Logger.mu.Unlock()
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	}
	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer
	serialized, err := entry.Logger.Formatter.Format(&entry)
	entry.Buffer = nil
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	} else {
		entry.Logger.mu.Lock()
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		entry.Logger.mu.Unlock()
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}
````

Common

- the Entry -> Formatter -> Output pipeline
- non pointer receiver on `Entry.log`
- support for field
  - gommon: never actually used, only the pkg field, `AddField` update in place
  - logrus: `WithField` create a new entry with added field

Differences with gommon/log v1

- has hook
- has `Writer()`, can be used for methods that accepts io.Writer
- has mutex on logger, lock/unlock when using hook
  - gommon: no hook thus no mutex
- has pool for bytes.Buffer, entry.Buffer is obtained from pool
- has pool for entry, `newEntry` and `releaseEntry`
  - gommon: each package has just one entry, created when init
- has json formatter
- has `Debug` etc. on `Logger`, implemented by create a new entry
  - gommon: user must create entry if they want to log
- has `*ln` functions, where `*` is `Debug`, `Info` etc.
- does not have concept of package
- does not have filter by package
- does not have source line