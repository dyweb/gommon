# 2018-12-02 Tree of logger

## Issues

- [#33](https://github.com/dyweb/gommon/issues/33) Maintain a Tree of Loggers is hard when struct has their own loggers
- [#78](https://github.com/dyweb/gommon/issues/78) Simplify the relationship of loggers

## Background

This is the original motivation of having gommon, have hierarchy between loggers and give application
fine grained control over its dependency libraries' log levels.
Most existing log libraries either requires library authors to expose logging configuration (i.e. pass a logger when
creating a client) or use a global logging level.
This result in hard time for developer & user to filter out useful message during development and when report errors.
A ideal scenario in pseudo code would be `server --log-debug=Server.WritePosts,Server.DeletePosts` where only some
package has debug level enabled

````text
// user report on already released code
// server --log-debug=Server.WritePosts,Server.DeletePosts
func (s *Server) WritePosts() {
   log := s.logger.Method() // a method logger
   log.Debug("very huge text")
}

// developer wrote a buggy code in a new method, hard code log level in code for the time being
func (s *Server) WritePosts() {
    // log := s.logger // normal logger
    log := s.logger.Copy().SetLevel(Debug)
    log.Debug("very huge text")
}
````

A good example would be the Java world, where you can use a xml file to config log of multiple packages.
[cihub/seelog](https://github.com/cihub/seelog) is a very close implementation in go.
Also the main inspiration was the log level config UI in solr

![solr-log-ui](../survey/solr-log-admin.png)

In summary the problem I wanted to solve by introducing a tree of logger is to
allow config log level of individual package/struct so that the log is filtered
from producer and no external filtering tool is required.

## Previous Solutions

The original solution in log v1 is a filter inspired by logrus's hook,
it keeps a set (actually it's a map) of package names, and entry struct
has a field called Pkg. Actually I think I never used it during the lifespan
of log v1, in v2 I started model after apex/log and uber-go/zap so hooks
are omitted

````go
// Filter determines if the entry should be logged
type Filter interface {
	Accept(entry *Entry) bool
	FilterName() string
	FilterDescription() string
}

var _ Filter = (*PkgFilter)(nil)

// PkgFilter only allows entry without `pkg` field or `pkg` value in the allow set to pass
// TODO: we should support level
// TODO: a more efficient way might be trie tree and use `/` to divide package into segments instead of using character
type PkgFilter struct {
	allow st.Set
}

// Accept checks if the entry.Pkg (NOT entry.Fields["pkg"]) is in the white list
func (filter *PkgFilter) Accept(entry *Entry) bool {
	return filter.allow.Contains(entry.Pkg)
}

// FilterName implements Filter interface
func (filter *PkgFilter) FilterName() string {
	return "PkgFilter"
}

func (filter *PkgFilter) FilterDescription() string {
	return "Filter log based on their pkg tag value, it is logged if it does not have pkg field or in whitelist"
}
````

Then follow the idea in the solr UI I decided to make **a tree of logger**
and allow each logger has their own logging level, because instead of
filtering based on package, you can simply turn one package to a more
verbose level while keep other package at their original level.
Thus I allow the loggers to contains children and identity,
parent children relationship is kept by forcing all the loggers created
using factory function.
Identity is obtained when create the logger using factory function using
runtime, I used some dirty trick to split method and function.
Then to set log level one have to traverse the logger tree to set all the
loggers, each package expose its top level logger, which is also a registry.

````go
type Logger struct {
	// mu is a Mutex instead of RWMutex because it's only for avoid concurrent write,
	// for performance reason and the natural of logging, reading stale config is not a big problem,
	// so we don't check mutex on read operation (i.e. log message) and allow race condition
	mu       sync.Mutex
	h        Handler
	level    Level
	source   bool
	children map[string][]*Logger

	id *Identity // use nil so we can have logger without identity
}

//
// TODO: allow release a child logger, this will be a trouble if we created 1,000 Client struct with its own logger...
func (l *Logger) AddChild(child *Logger) {
	l.mu.Lock()
	if l.children == nil {
		l.children = make(map[string][]*Logger, 1)
	}
	// children are group by their identity, i.e a package logger may have many struct logger of same struct because
	// that struct is used in multiple goroutines, those loggers have different pointer for identity, but they should
	// have same source line, so we use SourceLocation as key
	k := child.id.SourceLocation()
	children := l.children[k]
	// avoid putting same pointer twice, though it should never happen if AddChild is called correctly
	exists := false
	for _, c := range children {
		if c == child {
			exists = true
			break
		}
	}
	if !exists {
		l.children[k] = append(children, child)
	}
	l.mu.Unlock()
}

func SetLevelRecursive(root *Logger, level Level) {
	visited := make(map[*Logger]bool)
	PreOrderDFS(root, visited, func(l *Logger) {
		// TODO: remove it after we have tested it ....
		//fmt.Println(l.Identity().String())
		l.SetLevel(level)
	})
}

func NewPackageLogger() *Logger {
	return NewPackageLoggerWithSkip(1)
}

func NewFunctionLogger(packageLogger *Logger) *Logger {
	id := NewIdentityFromCaller(1)
	l := &Logger{
		id: &id,
	}
	return newLogger(packageLogger, l)
}

func NewStructLogger(packageLogger *Logger, loggable LoggableStruct) *Logger {
	id := loggable.LoggerIdentity(func() Identity {
		return NewIdentityFromCaller(1)
	})
	l := &Logger{
		id: &id,
	}
	l = newLogger(packageLogger, l)
	loggable.SetLogger(l)
	return l
}

func NewMethodLogger(structLogger *Logger) *Logger {
	id := NewIdentityFromCaller(1)
	l := &Logger{
		id: &id,
	}
	return newLogger(structLogger, l)
}
````

Keeping this parent child relationship introduce problem to gc because
after lifespan of a struct/method has expired, their logger is still referenced
by their parent logger in the map.

````go
var log = NewPackageLogger()

func handleEcho(w http.ResponseWriter, r *http.Request) {
  logger := NewFunctionLogger(log)
  logger.Info("handling echo")
}

func main() {
  m := http.NewServerMux()
  m.HandlerFunc("/echo", handleEcho)
  http.ListenAndServe(m)
}
````

In the example above, every time the echo func is called, a new entry is added
to the package level logger, the logger is useless after the func is finished
because next time it is going to create a new one, but because of
the children map in its parent, the logger will be referenced and can't be
garbage collected, not to mention have a huge map with pointer value will
also cause problem to garbage collector.

Besides the performance (it's actually memory leak)
We are also mixing the logger with logger registry but saving the map in registry.
One observation is package/library level logger are singleton, because
I just define them as `var log` in `pkg.go`, so they will be there in the entire
application lifecycle and there will be a fixed number of them.
Function and method logger will have many, it's a common practice to create
a new logger with fields attached based on function parameters, and they won't
run for a long time.
Struct is tricky, some struct like Server will likely have same span as entire
application. While others like Worker will only live for a short time.

````go
var log = NewPackageLogger()

func findUser(u string)  {
  logger := NewFunctionLogger(log).AddField(dlog.Str("user", u))
  logger.Info("start finding user")
  logger.Warn("user not found")
}

// long running struct like server
type Server struct {
  logger *dlog.Logger
}

func (s *Server) Echo(w http.ResponseWriter, r *http.Request)  {
  myName := r.Query().Get("name")
  logger := NewMethodLogger(s.logger).AddField(dlog.Str("name", myName))
  logger.Info("let's say hi")
}

// short running
type Worker struct {
  logger *dlog.Logger
}

func NewWorker() *Worker {
  w := Worker{}
  dlog.NewStructLogger(w) // need to generate getter setter using gommon
}

func (s *Worker) Fetch() {
  // wget google.com ....
}

func Work() {
  var wg sync.WaitGroup
  wg.Add(100)
  go func() {
    w := NewWoker()
    w.Fetch()
    wg.Done()
  }
  wg.Wait()
}
````

## Proposed solutions

Based on previous solutions we can find the following patterns

- keep the parent children relationship for all the loggers is unrealistic (gc) and unnecessary (a lot of them are short lived, just copy the parent's level and handler is fine)
- set level only gives some basic control, we can have more complex control based on logger location (the identity we have now) and caller location, this logic can be implemented in handler
- we can move logger relationship out into registry and only register important loggers like library, package and long running struct

Some performance concerns

- we will put identity caller into handler interface, identity and caller should be put as struct, 
pass pointer to struct may cause heap allocation, pass struct need to deal with empty identity ...

The new design for tree of logger and filter log when generating log has the following part

- registry that keeps a tree of registry and loggers
- handler that accept identity, user can implement any logic inside that handler

````text
type Logger struct {
    identity *Identity // could be nil, though most time it should not be except in benchmark and test
}

type Registry {
   childRegistry map[string]*Registry // use pointer instead of struct because Registry also use slice, if it is only using map, we can use struct because map is actually a pointer to underlying hashmap struct
   childLoggers []*Logger
}

func SetLevel(rootRegistry Registry, lvl Level) {
    DfsRegistry(rootRegistry, func(reg *Registry) {
        for _, l := range reg.Logger {
            l.SetLevel(lvl)
        }
    })
    // or
    DfsLogger(rootRegistry, func(l *Logger) {
        l.SetLevel(lvl)
    })
}

type Handler interface {
    HandleLog(loggerIdentity *Identity, caller Caller, level Level, msg string, context []Field, fields []Field)
}

type FileLineFilter struct {
    blockFile string
    blockLine int
    h Handler
}

func (fl *FileLineFilter) HandleLog(loggerIdentity *Identity, caller Caller, level Level, msg string, context []Field, fields []Field) {
    if caller.File == fl.blockFile {
        return
    }
    if caller.Line == fl.blockLine {
        return
    }
    // go through
    h(loggerIdentity, caller, msg, context, fields)
}

// TODO: need an example that check identity and level, if debug and is an allowed struct, go through
````


````text
// gommon/config/config.go
var log, logRegistry = NewPackageLoggerWithRegistry()

func NewConfigLoader() ConfigLoader {
    l := ConfigLoader()
    l.logger = dlog.NewStructLogger(log, l) // NOTE: package logger is now only used for copy level and handler, no longer used for registry
}

// ayi/web/server.go
var log, logRegistry = NewPackageLoggerWithRegistry()

func NewServer() Server {
    s := Server{}
    s.logger = dlog.NewStructLogger(log, s)
    logRegistry.AddLogger(s.logger) // server is long running and we know when it shuts down
}

func (s *Server) Echo(w http.ResponseWriter, r *http.Request) {
    logger := s.logger.MethodLogger() // add method field automatically
    
}
````

## Implementation

- [x] remove the parent children logic from logger
- [ ] add traverse registry and logger
  - [x] add traverse log in PreOrderDfs
- [ ] fix logic for func and method logger
- [ ] check if the skip caller is correct when create logger registry, unit test is in same package so it's always correct
- [ ] check if identity of log registry is correct
- [ ] check if identity of logger is correct
- [x] change caller to struct
  - [ ] it may have some performance impact