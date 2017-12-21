# Apex-Log

https://github.com/apex/log

- complexity: medium
- use handler instead of formatter + writer
  - like `net/http` you can use a function as handler
  - [ ] TODO: why it is `func (f HandlerFunc)` instead of `func (f *HandlerFunc)`, can we have a `type HandlerFuncPtr`?

````go
func AuthWrapper(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // some auth logic ...
        h.ServeHTTP(w, r)
    })
}
````

entry -> formatter -> bytes -> writer
entry -> handler (convert entry to format accepted by destination)

Entry

````go
// Entry represents a single log entry.
type Entry struct {
	Logger    *Logger   `json:"-"`
	Fields    Fields    `json:"fields"`
	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	start     time.Time
	fields    []Fields
}


// WithFields returns a new entry with `fields` set.
func (e *Entry) WithFields(fields Fielder) *Entry {
	f := []Fields{}
	f = append(f, e.fields...)
	f = append(f, fields.Fields())
	return &Entry{
		Logger: e.Logger,
		fields: f,
	}
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Logger.log(DebugLevel, e, msg)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string) {
	e.Logger.log(FatalLevel, e, msg)
	os.Exit(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string) *Entry {
	return &Entry{
		Logger:    e.Logger,
		Fields:    e.mergedFields(),
		Level:     level,
		Message:   msg,
		Timestamp: Now(),
	}
}
````

Logger

- [x] TODO: the comment on `log` seems to be stale, there is no clone and the check of level is done before passing to handler
  - clone is removed in [refactor Entry using slice of Fields to improve perf](https://github.com/apex/log/commit/ae8aa5030551bd783ae9c82f281c50c2ec32ba4f)

````go
// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	HandleLog(*Entry) error
}

// Logger represents a logger with configurable Level and Handler.
type Logger struct {
	Handler Handler
	Level   Level
}

// WithFields returns a new entry with `fields` set.
func (l *Logger) WithFields(fields Fielder) *Entry {
	return NewEntry(l).WithFields(fields.Fields())
}

// Debug level message.
func (l *Logger) Debug(msg string) {
	NewEntry(l).Debug(msg)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not
// met.
func (l *Logger) log(level Level, e *Entry, msg string) {
	if level < l.Level {
		return
	}

	if err := l.Handler.HandleLog(e.finalize(level, msg)); err != nil {
		stdlog.Printf("error logging: %s", err)
	}
}
````

Handler

- handlers use lock, though I don't know if it is really necessary, 
maybe it is only needed when you call `fmt.Fprintf` several times to log a single entry

````go
// cli.go
// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "\033[%dm%*s\033[0m %-25s", color, h.Padding+1, level, e.Message)

	for _, name := range names {
		if name == "source" {
			continue
		}

		fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
````

````go
// multi.go
// Handler implementation.
type Handler struct {
	Handlers []log.Handler
}

// New handler.
func New(h ...log.Handler) *Handler {
	return &Handler{
		Handlers: h,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	for _, handler := range h.Handlers {
		// TODO(tj): maybe just write to stderr here, definitely not ideal
		// to miss out logging to a more critical handler if something
		// goes wrong
		if err := handler.HandleLog(e); err != nil {
			return err
		}
	}

	return nil
}
````