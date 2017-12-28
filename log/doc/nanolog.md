# nanlog

https://github.com/ScottMansfield/nanolog

- complexity: medium
- use binary format, don't store the unchanged part (log format) multiple times
  - more like a database file with dictionary encoding

````go
var h nanolog.Handle

func init() {
	nanolog.SetWriter(os.Stderr)
	h = nanolog.AddLogger("Example %i32 log %{s} line %c128")
}

func main() {
	nanolog.Log(h, int32(4), "this is a string", 4i)
	nanolog.Flush()
}
````

- when `AddLogger` gets log format, parse it and store the expected types in `Logger.Kinds []reflect.Kind`
- when `Log`, use `reflect.TypeOf(args[idx]).Kind()`, use `binary` package to write to buf then flush to file

````go
// Handle is a simple handle to an internal logging data structure
// LogHandles are returned by the AddLogger method and used by the Log method to
// actually log data.
type Handle uint32

// Logger is the internal struct representing the runtime state of the loggers.
// The Segs field is not used during logging; it is only used in the inflate
// utility
type Logger struct {
	Kinds []reflect.Kind
	Segs  []string
}

type logWriter struct {
	initBuf  *bytes.Buffer
	w        *bufio.Writer
	firstSet bool

	writeLock sync.Locker

	loggers       []Logger
	curLoggersIdx *uint32
}

func (lw *logWriter) AddLogger(fmt string) Handle {
	// save some kind of string format to the file
	idx := atomic.AddUint32(lw.curLoggersIdx, 1) - 1
	l, segs := parseLogLine(fmt)
	lw.loggers[idx] = l
	lw.writeLogDataToFile(idx, l.Kinds, segs)
	return Handle(idx)
}

func (lw *logWriter) Log(handle Handle, args ...interface{}) error {
	l := lw.loggers[handle]

	if len(l.Kinds) != len(args) {
		panic("Number of args does not match log line")
	}

	buf := bufpool.Get().(*[]byte)
	*buf = (*buf)[:0]
	b := make([]byte, 8)

	*buf = append(*buf, byte(ETLogEntry))

	binary.LittleEndian.PutUint32(b, uint32(handle))
	*buf = append(*buf, b[:4]...)

	for idx := range l.Kinds {
		if l.Kinds[idx] != reflect.TypeOf(args[idx]).Kind() {
			panic("Argument type does not match log line")
		}

		// write serialized version to writer
		switch l.Kinds[idx] {
		case reflect.Bool:
			if args[idx].(bool) {
				*buf = append(*buf, 1)
			} else {
				*buf = append(*buf, 0)
			}

		case reflect.String:
			s := args[idx].(string)
			binary.LittleEndian.PutUint32(b, uint32(len(s)))
			*buf = append(*buf, b[:4]...)
			*buf = append(*buf, s...)

		// ints
		case reflect.Int:
			// Assume generic int is 64 bit
			i := args[idx].(int)
			binary.LittleEndian.PutUint64(b, uint64(i))
			*buf = append(*buf, b...)

		// floats
		case reflect.Float32:
			f := args[idx].(float32)
			i := math.Float32bits(f)
			binary.LittleEndian.PutUint32(b, i)
			*buf = append(*buf, b[:4]...)

		default:
			panic(fmt.Sprintf("Invalid Kind in logger: %v", l.Kinds[idx]))
		}
	}

	lw.writeLock.Lock()
	_, err := lw.w.Write(*buf)
	lw.writeLock.Unlock()

	bufpool.Put(buf)
	return err
}
````