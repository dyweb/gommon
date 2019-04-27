package log

import (
	"sort"
	"sync"

	"github.com/dyweb/gommon/util/runtimeutil"
)

// registry.go is used for maintain relationship between loggers across packages and projects
// it also contains util func for traverse registry and logger

var globalRegistryGroup = newRegistryGroup()

type registryGroup struct {
	// we didn't use RWMutex because when walking a group, the main purpose is to modify loggers
	// inside registries, so two walking should not happens in parallel
	mu         sync.Mutex
	registries map[string]*Registry
}

func newRegistryGroup() *registryGroup {
	return &registryGroup{
		registries: make(map[string]*Registry),
	}
}

func (rg *registryGroup) add(reg *Registry) {
	rg.mu.Lock()
	defer rg.mu.Unlock()

	id := reg.identity
	if id == "" {
		panic("log registry identity is empty")
	}
	oReg, ok := rg.registries[id]
	if ok {
		if oReg == reg {
			return
		} else {
			panic("log registry is already registered for " + id)
		}
	}
	rg.registries[id] = reg
}

// Registry contains default and tracked loggers, it is per package
type Registry struct {
	mu      sync.Mutex
	loggers []*Logger

	// identity is a string for package
	identity string
}

// NewRegistry create a log registry with a default logger for a package.
// It registers itself in globalRegistryGroup so it can be updated later using WalkRegistries
func NewRegistry() *Registry {
	frame := runtimeutil.GetCallerFrame(1)
	pkg, _ := runtimeutil.SplitPackageFunc(frame.Function)
	reg := Registry{
		identity: pkg,
		loggers:  []*Logger{NewPackageLoggerWithSkip(1)},
	}
	globalRegistryGroup.add(&reg)
	return &reg
}

func (r *Registry) Identity() string {
	return r.identity
}

// Logger returns the default logger in registry
func (r *Registry) Logger() *Logger {
	if len(r.loggers) < 1 {
		panic("no default logger found in registry")
	}
	return r.loggers[0]
}

// NewLogger creates a logger based on default logger and register it in registry.
// It should be used sparingly, if you need to add more fields as context for a func,
// you should make copy from default logger using methods like TODO: WithFields?
// to avoid register them in registry
func (r *Registry) NewLogger() *Logger {
	// no lock is added because addLogger also acquire lock
	id := NewIdentityFromCaller(1)
	l := Logger{
		id: &id,
	}
	newLogger(r.Logger(), &l)
	r.addLogger(&l)
	return &l
}

// addLogger registers a logger into registry. It's a nop if the logger is already there
func (r *Registry) addLogger(l *Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ol := range r.loggers {
		if ol == l {
			return
		}
	}
	r.loggers = append(r.loggers, l)
}

// WalkRegistry walks registry in globalRegistryGroup in sorted order of id (package path)
func WalkRegistry(cb func(r *Registry)) {
	group := globalRegistryGroup
	group.mu.Lock()
	defer group.mu.Unlock()

	// visit in the order of sorted id
	var ids []string
	for id := range group.registries {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		cb(group.registries[id])
	}
}

// WalkLogger calls WalkRegistry and within each registry, walk in insert order of loggers
func WalkLogger(cb func(l *Logger)) {
	WalkRegistry(func(r *Registry) {
		r.mu.Lock()
		defer r.mu.Unlock()

		for _, l := range r.loggers {
			cb(l)
		}
	})
}
