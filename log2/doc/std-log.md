# Log

https://golang.org/src/log/log.go

- it locks the logger when output
- output is `io.Writer`
- the buffer seems to be shared globally, that's why the unlock is using `defer`
- getting caller info is expensive ... that's why it should be done at compile time ...
  - Java need to use runtime information, get stack to get file and line number as well
- std logger goes to standard error

````go
var std = New(os.Stderr, "", LstdFlags)

func Printf(format string, v ...interface{}) {
  	std.Output(2, fmt.Sprintf(format, v...))
}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
  	mu     sync.Mutex // ensures atomic writes; protects the following fields
  	prefix string     // prefix to write at beginning of each line
  	flag   int        // properties
  	out    io.Writer  // destination for output
  	buf    []byte     // for accumulating text to write
}
  
func (l *Logger) Printf(format string, v ...interface{}) {
  	l.Output(2, fmt.Sprintf(format, v...))
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
  	// Get time early if we need it.
  	var now time.Time
  	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
  		now = time.Now()
  	}
  	var file string
  	var line int
  	l.mu.Lock()
  	defer l.mu.Unlock()
  	if l.flag&(Lshortfile|Llongfile) != 0 {
  		// Release lock while getting caller info - it's expensive.
  		l.mu.Unlock()
  		var ok bool
  		_, file, line, ok = runtime.Caller(calldepth)
  		if !ok {
  			file = "???"
  			line = 0
  		}
  		l.mu.Lock()
  	}
  	l.buf = l.buf[:0]
  	l.formatHeader(&l.buf, now, file, line)
  	l.buf = append(l.buf, s...)
  	if len(s) == 0 || s[len(s)-1] != '\n' {
  		l.buf = append(l.buf, '\n')
  	}
  	_, err := l.out.Write(l.buf)
  	return err
}
````