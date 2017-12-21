# Log

https://golang.org/src/log/log.go

- complexity: low
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

Printf src/fmt/print.go, I think compared with zap, reflection is not the problem, it does not use reflection for simple types,
and zap is using switch as well, except based on pre stored int (represent type) instead of `v.(type)`

- https://blog.golang.org/laws-of-reflection
- I remember I seems to have played with reflection when trying to add strcut tag support for [go-solr#11](https://github.com/at15/go-solr/issues/11)

````go
func (p *pp) doPrintf(format string, a []interface{}) {
	// .... 
    p.printArg(a[argNum], verb)
    argNum++
}

func (p *pp) printArg(arg interface{}, verb rune) {
	p.arg = arg
	p.value = reflect.Value{}

	if arg == nil {
		switch verb {
		case 'T', 'v':
			p.fmt.padString(nilAngleString)
		default:
			p.badVerb(verb)
		}
		return
	}

	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb {
	case 'T':
		p.fmt.fmt_s(reflect.TypeOf(arg).String())
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), 'p')
		return
	}

	// Some types can be done without reflection.
	switch f := arg.(type) {
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
		p.fmtFloat(f, 64, verb)
	case complex64:
		p.fmtComplex(complex128(f), 64, verb)
	case complex128:
		p.fmtComplex(f, 128, verb)
	case int:
		p.fmtInteger(uint64(f), signed, verb)
	case int8:
		p.fmtInteger(uint64(f), signed, verb)
	case int16:
		p.fmtInteger(uint64(f), signed, verb)
	case int32:
		p.fmtInteger(uint64(f), signed, verb)
	case int64:
		p.fmtInteger(uint64(f), signed, verb)
	case uint:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint8:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint16:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint32:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint64:
		p.fmtInteger(f, unsigned, verb)
	case uintptr:
		p.fmtInteger(uint64(f), unsigned, verb)
	case string:
		p.fmtString(f, verb)
	case []byte:
		p.fmtBytes(f, verb, "[]byte")
	case reflect.Value:
		// Handle extractable values with special methods
		// since printValue does not handle them at depth 0.
		if f.IsValid() && f.CanInterface() {
			p.arg = f.Interface()
			if p.handleMethods(verb) {
				return
			}
		}
		p.printValue(f, verb, 0)
	default:
		// If the type is not simple, it might have methods.
		if !p.handleMethods(verb) {
			// Need to use reflection, since the type had no
			// interface methods that could be used for formatting.
			p.printValue(reflect.ValueOf(f), verb, 0)
		}
	}
}


// printValue is similar to printArg but starts with a reflect value, not an interface{} value.
// It does not handle 'p' and 'T' verbs because these should have been already handled by printArg.
func (p *pp) printValue(value reflect.Value, verb rune, depth int) {
	// Handle values with special methods if not already handled by printArg (depth == 0).
	if depth > 0 && value.IsValid() && value.CanInterface() {
		p.arg = value.Interface()
		if p.handleMethods(verb) {
			return
		}
	}
	p.arg = nil
	p.value = value

	switch f := value; value.Kind() {
	case reflect.Invalid:
		if depth == 0 {
			p.buf.WriteString(invReflectString)
		} else {
			switch verb {
			case 'v':
				p.buf.WriteString(nilAngleString)
			default:
				p.badVerb(verb)
			}
		}
	case reflect.Bool:
		p.fmtBool(f.Bool(), verb)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.fmtInteger(uint64(f.Int()), signed, verb)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.fmtInteger(f.Uint(), unsigned, verb)
	case reflect.Float32:
		p.fmtFloat(f.Float(), 32, verb)
	case reflect.Float64:
		p.fmtFloat(f.Float(), 64, verb)
	case reflect.Complex64:
		p.fmtComplex(f.Complex(), 64, verb)
	case reflect.Complex128:
		p.fmtComplex(f.Complex(), 128, verb)
	case reflect.String:
		p.fmtString(f.String(), verb)
	case reflect.Map:
		if p.fmt.sharpV {
			p.buf.WriteString(f.Type().String())
			if f.IsNil() {
				p.buf.WriteString(nilParenString)
				return
			}
			p.buf.WriteByte('{')
		} else {
			p.buf.WriteString(mapString)
		}
		keys := f.MapKeys()
		for i, key := range keys {
			if i > 0 {
				if p.fmt.sharpV {
					p.buf.WriteString(commaSpaceString)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			p.printValue(key, verb, depth+1)
			p.buf.WriteByte(':')
			p.printValue(f.MapIndex(key), verb, depth+1)
		}
		if p.fmt.sharpV {
			p.buf.WriteByte('}')
		} else {
			p.buf.WriteByte(']')
		}
	case reflect.Struct:
		if p.fmt.sharpV {
			p.buf.WriteString(f.Type().String())
		}
		p.buf.WriteByte('{')
		for i := 0; i < f.NumField(); i++ {
			if i > 0 {
				if p.fmt.sharpV {
					p.buf.WriteString(commaSpaceString)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			if p.fmt.plusV || p.fmt.sharpV {
				if name := f.Type().Field(i).Name; name != "" {
					p.buf.WriteString(name)
					p.buf.WriteByte(':')
				}
			}
			p.printValue(getField(f, i), verb, depth+1)
		}
		p.buf.WriteByte('}')
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			if p.fmt.sharpV {
				p.buf.WriteString(f.Type().String())
				p.buf.WriteString(nilParenString)
			} else {
				p.buf.WriteString(nilAngleString)
			}
		} else {
			p.printValue(value, verb, depth+1)
		}
	case reflect.Array, reflect.Slice:
		switch verb {
		case 's', 'q', 'x', 'X':
			// Handle byte and uint8 slices and arrays special for the above verbs.
			t := f.Type()
			if t.Elem().Kind() == reflect.Uint8 {
				var bytes []byte
				if f.Kind() == reflect.Slice {
					bytes = f.Bytes()
				} else if f.CanAddr() {
					bytes = f.Slice(0, f.Len()).Bytes()
				} else {
					// We have an array, but we cannot Slice() a non-addressable array,
					// so we build a slice by hand. This is a rare case but it would be nice
					// if reflection could help a little more.
					bytes = make([]byte, f.Len())
					for i := range bytes {
						bytes[i] = byte(f.Index(i).Uint())
					}
				}
				p.fmtBytes(bytes, verb, t.String())
				return
			}
		}
		if p.fmt.sharpV {
			p.buf.WriteString(f.Type().String())
			if f.Kind() == reflect.Slice && f.IsNil() {
				p.buf.WriteString(nilParenString)
				return
			}
			p.buf.WriteByte('{')
			for i := 0; i < f.Len(); i++ {
				if i > 0 {
					p.buf.WriteString(commaSpaceString)
				}
				p.printValue(f.Index(i), verb, depth+1)
			}
			p.buf.WriteByte('}')
		} else {
			p.buf.WriteByte('[')
			for i := 0; i < f.Len(); i++ {
				if i > 0 {
					p.buf.WriteByte(' ')
				}
				p.printValue(f.Index(i), verb, depth+1)
			}
			p.buf.WriteByte(']')
		}
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if depth == 0 && f.Pointer() != 0 {
			switch a := f.Elem(); a.Kind() {
			case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
				p.buf.WriteByte('&')
				p.printValue(a, verb, depth+1)
				return
			}
		}
		fallthrough
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		p.fmtPointer(f, verb)
	default:
		p.unknownType(f)
	}
}
````