# Zap

https://github.com/uber-go/zap

Zap focus on speed, reduce allocation

- complexity: high
- standard logger use `fmt.Sprintf`, it is slow because it is based on reflection, also apply to `encoding/json`
- it has `SugaredLogger` with a `Printf` like API, it is based on Logger
- Logger requires using `zapcore.Field` in `fields.go`

````go
logger, _ := zap.NewProduction()
defer logger.Sync()
logger.Info("failed to fetch URL",
  // Structured context as strongly typed Field values.
  zap.String("url", url),
  zap.Int("attempt", 3),
  zap.Duration("backoff", time.Second),
)
````

Logger

- `Check`, is like `isDebug` in log4j
  - calls `check`, which call `log.core.Check`, the core will add itself to the checked entry
- `Debug` (along w/ `Info` etc.)
  - call `check` to get `zapcore.CheckedEntry`
    - call `log.core.Check`, core will check level and add itself to `CheckedEntry.cores` using `AddCore`
      - `AddCore` call `getCheckedEntry` to get entry from global pool if receiver is `nil`
  - call `CheckedEntry.Write(fields ...Field)` to write fields (context) to `cores`
    - call `ce.cores[i].Write(ce.Entry, fields)`
      - `Core.Write` call `c.enc.EncodeEntry` (i `Encoder`) to get a buffer and write to `c.out` (i `WriteSyncer`)
    - it will check `should` to decide if panic of fatal is needed ...
- encoding logic can be find in `zapcore/console_encoder.go`, `zapcore/field.go`, `zapcore/json_encoder.go`
  - which eventually goes to `buffer/buffer.go` and call like `strconv.AppendInt(b.bs, i, 10)` where `bs` is `[]byte`
  - **Unlike the standard library's bytes.Buffer, it supports a portion of the strconv package's zero-allocation formatters**
  - `(f Field) AddTo` is used to `switch` based on `f.Type` and call `AddArray`, `AddObject` etc. on encoder
  - `AddReflected` is used if a field is created without concrete type, and use `json.Marshal`
  
Extra

- sampler, it wraps a existing core and check if log should be sampled
- tee, duplicate one core to multiple cores
  - [ ] TODO: then why CheckedEntry has `cores` instead of just one core ... seems never used
- memory_encoder.go used for testing, use `map[string]interface{}`
- console_encoder.go, encode fields as json, but message and timestamp as plain text

````go
// in logger.go
type Logger struct {
	core zapcore.Core

	development bool
	name        string
	errorOutput zapcore.WriteSyncer

	addCaller bool
	addStack  zapcore.LevelEnabler

	callerSkip int
}

type SugaredLogger struct {
	base *Logger
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func (log *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return log.check(lvl, msg)
}

func (log *Logger) check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	// check must always be called directly by a method in the Logger interface
	// (e.g., Check, Info, Fatal).
	const callerSkipOffset = 2

	// Create basic checked entry thru the core; this will be non-nil if the
	// log message will actually be written somewhere.
	ent := zapcore.Entry{
		LoggerName: log.name,
		Time:       time.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := log.core.Check(ent, nil)
	willWrite := ce != nil

	// Set up any required terminal behavior.
	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.Should(ent, zapcore.WriteThenPanic)
	case zapcore.FatalLevel:
		ce = ce.Should(ent, zapcore.WriteThenFatal)
	case zapcore.DPanicLevel:
		if log.development {
			ce = ce.Should(ent, zapcore.WriteThenPanic)
		}
	}

	// Only do further annotation if we're going to write this message; checked
	// entries that exist only for terminal behavior don't benefit from
	// annotation.
	if !willWrite {
		return ce
	}

	// Thread the error output through to the CheckedEntry.
	ce.ErrorOutput = log.errorOutput
	if log.addCaller {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(log.callerSkip + callerSkipOffset))
		if !ce.Entry.Caller.Defined {
			fmt.Fprintf(log.errorOutput, "%v Logger.check error: failed to get caller\n", time.Now().UTC())
			log.errorOutput.Sync()
		}
	}
	if log.addStack.Enabled(ce.Entry.Level) {
		ce.Entry.Stack = Stack("").String
	}

	return ce
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Debug(msg string, fields ...zapcore.Field) {
	if ce := log.check(DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}
````

Entry & CheckedEntry entry.go

````go
// An Entry represents a complete log message. The entry's structured context
// is already serialized, but the log level, time, message, and call site
// information are available for inspection and modification.
//
// Entries are pooled, so any functions that accept them MUST be careful not to
// retain references to them.
type Entry struct {
	Level      Level
	Time       time.Time
	LoggerName string
	Message    string
	Caller     EntryCaller
	Stack      string
}

// CheckedEntry is an Entry together with a collection of Cores that have
// already agreed to log it.
//
// CheckedEntry references should be created by calling AddCore or Should on a
// nil *CheckedEntry. References are returned to a pool after Write, and MUST
// NOT be retained after calling their Write method.
type CheckedEntry struct {
	Entry
	ErrorOutput WriteSyncer
	dirty       bool // best-effort detection of pool misuse
	should      CheckWriteAction
	cores       []Core
}

// AddCore adds a Core that has agreed to log this CheckedEntry. It's intended to be
// used by Core.Check implementations, and is safe to call on nil CheckedEntry
// references.
func (ce *CheckedEntry) AddCore(ent Entry, core Core) *CheckedEntry {
	if ce == nil {
		ce = getCheckedEntry()
		ce.Entry = ent
	}
	ce.cores = append(ce.cores, core)
	return ce
}

// Write writes the entry to the stored Cores, returns any errors, and returns
// the CheckedEntry reference to a pool for immediate re-use. Finally, it
// executes any required CheckWriteAction.
func (ce *CheckedEntry) Write(fields ...Field) {
	if ce == nil {
		return
	}

	if ce.dirty {
		if ce.ErrorOutput != nil {
			// Make a best effort to detect unsafe re-use of this CheckedEntry.
			// If the entry is dirty, log an internal error; because the
			// CheckedEntry is being used after it was returned to the pool,
			// the message may be an amalgamation from multiple call sites.
			fmt.Fprintf(ce.ErrorOutput, "%v Unsafe CheckedEntry re-use near Entry %+v.\n", time.Now(), ce.Entry)
			ce.ErrorOutput.Sync()
		}
		return
	}
	ce.dirty = true

	var err error
	for i := range ce.cores {
		err = multierr.Append(err, ce.cores[i].Write(ce.Entry, fields))
	}
	if ce.ErrorOutput != nil {
		if err != nil {
			fmt.Fprintf(ce.ErrorOutput, "%v write error: %v\n", time.Now(), err)
			ce.ErrorOutput.Sync()
		}
	}

	should, msg := ce.should, ce.Message
	putCheckedEntry(ce)

	switch should {
	case WriteThenPanic:
		panic(msg)
	case WriteThenFatal:
		exit.Exit()
	}
}

var (
	_cePool = sync.Pool{New: func() interface{} {
		// Pre-allocate some space for cores.
		return &CheckedEntry{
			cores: make([]Core, 4),
		}
	}}
)

func getCheckedEntry() *CheckedEntry {
	ce := _cePool.Get().(*CheckedEntry)
	ce.reset()
	return ce
}

func putCheckedEntry(ce *CheckedEntry) {
	if ce == nil {
		return
	}
	_cePool.Put(ce)
}
````

Field field.go `Any(key string, value interface{}) zapcore.Field` and `zapcore/field.go` `func (f Field) AddTo(enc ObjectEncoder)`

````go
// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) zapcore.Field {
	switch val := value.(type) {
	case zapcore.ObjectMarshaler:
		return Object(key, val)
	case zapcore.ArrayMarshaler:
		return Array(key, val)
	case bool:
		return Bool(key, val)
	case []bool:
		return Bools(key, val)
	case error:
		return NamedError(key, val)
	case []error:
		return Errors(key, val)
	case fmt.Stringer:
		return Stringer(key, val)
	default:
		return Reflect(key, val)
	}
}


// AddTo exports a field through the ObjectEncoder interface. It's primarily
// useful to library authors, and shouldn't be necessary in most applications.
func (f Field) AddTo(enc ObjectEncoder) {
	var err error

	switch f.Type {
	case ArrayMarshalerType:
		err = enc.AddArray(f.Key, f.Interface.(ArrayMarshaler))
	case ObjectMarshalerType:
		err = enc.AddObject(f.Key, f.Interface.(ObjectMarshaler))
	case ReflectType:
		err = enc.AddReflected(f.Key, f.Interface)
	case NamespaceType:
		enc.OpenNamespace(f.Key)
	case StringerType:
		enc.AddString(f.Key, f.Interface.(fmt.Stringer).String())
	case ErrorType:
		encodeError(f.Key, f.Interface.(error), enc)
	case SkipType:
		break
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}

	if err != nil {
		enc.AddString(fmt.Sprintf("%sError", f.Key), err.Error())
	}
}
````

Core zapcore/core.go

````go
// Core is a minimal, fast logger interface. It's designed for library authors
// to wrap in a more user-friendly API.
type Core interface {
	LevelEnabler

	// With adds structured context to the Core.
	With([]Field) Core
	// Check determines whether the supplied Entry should be logged (using the
	// embedded LevelEnabler and possibly some extra logic). If the entry
	// should be logged, the Core adds itself to the CheckedEntry and returns
	// the result.
	//
	// Callers must use Check before calling Write.
	Check(Entry, *CheckedEntry) *CheckedEntry
	// Write serializes the Entry and any Fields supplied at the log site and
	// writes them to their destination.
	//
	// If called, Write should always log the Entry and Fields; it should not
	// replicate the logic of Check.
	Write(Entry, []Field) error
	// Sync flushes buffered logs (if any).
	Sync() error
}

type ioCore struct {
	LevelEnabler
	enc Encoder
	out WriteSyncer
}

func (c *ioCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *ioCore) Write(ent Entry, fields []Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		c.Sync()
	}
	return nil
}
````

Encoder zapcore/encoder.go

````go
// Encoder is a format-agnostic interface for all log entry marshalers. Since
// log encoders don't need to support the same wide range of use cases as
// general-purpose marshalers, it's possible to make them faster and
// lower-allocation.
//
// Implementations of the ObjectEncoder interface's methods can, of course,
// freely modify the receiver. However, the Clone and EncodeEntry methods will
// be called concurrently and shouldn't modify the receiver.
type Encoder interface {
	ObjectEncoder

	// Clone copies the encoder, ensuring that adding fields to the copy doesn't
	// affect the original.
	Clone() Encoder

	// EncodeEntry encodes an entry and fields, along with any accumulated
	// context, into a byte buffer and returns it.
	EncodeEntry(Entry, []Field) (*buffer.Buffer, error)
}
````

Pool & Buffer buffer/buffer.go & pool.go

````go
// Package buffer provides a thin wrapper around a byte slice. Unlike the
// standard library's bytes.Buffer, it supports a portion of the strconv
// package's zero-allocation formatters.
package buffer

import (
	"strconv"
	"sync")

const _size = 1024 // by default, create 1 KiB buffers

// Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// the only way to construct one is via a Pool.
type Buffer struct {
	bs   []byte
	pool Pool
}

// A Pool is a type-safe wrapper around a sync.Pool.
type Pool struct {
	p *sync.Pool
}
````