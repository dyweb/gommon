# 2018-11-24 Fields and Handler interface

Follow up on the [survey](2018-11-18-design-continued.md) after trying out how other library support adding fields to 'logger'.

There are two ways, create a new entry that hold reference to the original logger or create a new logger.
I don't like the entry way, in this way the entry becomes the real logger, because all the user facing log methods like
`Debug`, `Info` all lands on entry instead of the logger. 
For the logger way, you need to clone the logger so the fields are added instead of updated, however for zap and zerolog,
in order to increase speed, they encode the new fields into bytes right away. zerolog is actually more entry way,
it is using `Event`, event copy logger's context (which is bytes) and writer (basically io.Writer) 

Actually the `Event` in zerolog is still different from `Entry` in logrus, in zerolog you are keep updating the context
of current event when you add fields, while in logrus you are creating new entry when adding fields.

````go
// apex/log/entry.go
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

// apex/log/entry.go
func (e *Entry) Debug(msg string) {
	e.Logger.log(DebugLevel, e, msg)
}
````

````go
// logrus/entry.go
func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}
````

````go
// zap/logger.go

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (log *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return log
	}
	l := log.clone()
	l.core = l.core.With(fields)
	return l
}
````

````go
// zerolog/event.go

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
type Event struct {
	buf   []byte
	w     LevelWriter
	level Level
	done  func(msg string)
	ch    []Hook // hooks from context
}


// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return e
	}
	e.buf = enc.AppendString(enc.AppendKey(e.buf, key), val)
	return e
}
````

## Proposed change set

- enable fields inside logger struct to keep context
  - `AddField(f Field)` to add field in place, normally used for struct and method logger
  - ~~`Entry()` to return a log entry that similar to entry in logrus & event in zerolog~~
- ~~`Entry` copies handler, config from the logger when it is created and never update, it should be short lived~~
  - `AddField` 
  - only have `Info` `Info` but no `InfoF` like logger
  - I don't think we really need entry when we have InfoF ... 

TODO

- [x] ~~change identity from pointer to struct~~ still keep as pointer so we can have nil value
  - for function that returns identity, return value instead of pointer
  - [ ] might change it to a string?
- [ ] enable fields
- [ ] allow add fields (don't allow remove)
- [ ] change handler interface or logger to merge fields with adhoc fields
  - I don't want to create new slice just for merging the fields, and logger need to check if its context is empty


- now it becomes permutation
  - source
  - fields
  - context
- we have total 2 * 2 * 2 ... 8 methods to implement, well ... such copy and paste and if else ...

````go
type Handler interface {
	HandleLogWithFields(level Level, time time.Time, msg string, fields Fields)
	// context are fields from the logger
	HandleLogWithContextFields(lvl Level, time time.Time, msg string, context Fields, fields Fields)
}
````