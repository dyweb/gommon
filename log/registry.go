package log

import (
	"sync"

	"github.com/dyweb/gommon/util/runtimeutil"
)

// registry.go is used for maintain relationship between loggers across packages and projects
// it also contains util func for traverse registry and logger

// Registry contains child registry and loggers
type Registry struct {
	mu       sync.Mutex
	children []*Registry
	loggers  []*Logger

	// immutable
	identity RegistryIdentity
}

type RegistryType uint8

const (
	UnknownRegistry RegistryType = iota
	ApplicationRegistry
	LibraryRegistry
	PackageRegistry
)

func (r RegistryType) String() string {
	switch r {
	case UnknownRegistry:
		return "unk"
	case ApplicationRegistry:
		return "app"
	case LibraryRegistry:
		return "lib"
	case PackageRegistry:
		return "pkg"
	default:
		return "unk"
	}
}

type RegistryIdentity struct {
	// Project is specified by user, i.e. for all the packages under gommon, they would have github.com/dyweb/gommon
	Project string
	// Package is detected base on runtime, i.e. github.com/dyweb/gommon/noodle
	Package string
	// Type is specified by user when creating registry
	Type RegistryType
	// File is where create registry is called
	File string
	// Line is where create registry is called
	Line int
}

func NewLibraryRegistry(project string) Registry {
	return Registry{
		identity: newRegistryId(project, LibraryRegistry, 0),
	}
}

// TODO: validate skip
func NewApplicationLoggerAndRegistry(project string) (*Logger, *Registry) {
	reg := Registry{
		identity: newRegistryId(project, ApplicationRegistry, 1),
	}
	logger := NewPackageLoggerWithSkip(1)
	reg.AddLogger(logger)
	return logger, &reg
}

func NewPackageLoggerAndRegistryWithSkip(project string, skip int) (*Logger, *Registry) {
	reg := Registry{
		identity: newRegistryId(project, PackageRegistry, skip+1),
	}
	logger := NewPackageLoggerWithSkip(skip + 1)
	reg.AddLogger(logger)
	return logger, &reg
}

func newRegistryId(proj string, tpe RegistryType, skip int) RegistryIdentity {
	// TODO: check if the skip works .... we need another package for testing that
	frame := runtimeutil.GetCallerFrame(skip + 1)
	pkg, _ := runtimeutil.SplitPackageFunc(frame.Function)
	return RegistryIdentity{
		Project: proj,
		Package: pkg,
		Type:    tpe,
		File:    frame.File,
		Line:    frame.Line,
	}
}

// AddRegistry is for adding a package level log registry to a library/application level log registry
// It skips add if child registry already there
func (r *Registry) AddRegistry(child *Registry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.children {
		if c == child {
			return
		}
	}
	r.children = append(r.children, child)
}

// AddLogger is used for registering a logger to package level log registry
// It skips add if the logger is already there
func (r *Registry) AddLogger(l *Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ol := range r.loggers {
		if ol == l {
			return
		}
	}
	r.loggers = append(r.loggers, l)
}

func (r *Registry) Identity() RegistryIdentity {
	return r.identity
}

func SetLevel(root *Registry, level Level) {
	WalkLogger(root, func(l *Logger) {
		l.SetLevel(level)
	})
}

func SetHandler(root *Registry, handler Handler) {
	WalkLogger(root, func(l *Logger) {
		l.SetHandler(handler)
	})
}

func EnableSource(root *Registry) {
	WalkLogger(root, func(l *Logger) {
		l.EnableSource()
	})
}

func DisableSource(root *Registry) {
	WalkLogger(root, func(l *Logger) {
		l.DisableSource()
	})
}

// WalkLogger is PreOrderDfs
func WalkLogger(root *Registry, cb func(l *Logger)) {
	walkLogger(root, nil, nil, cb)
}

// walkLogger loops loggers in current registry first, then visist its children in DFS
// It will create the map if it is nil, so caller don't need to do the bootstrap,
// However, caller can provide a map so they can use it to get all the loggers and registry
func walkLogger(root *Registry, loggers map[*Logger]bool, registries map[*Registry]bool, cb func(l *Logger)) {
	// first call
	if loggers == nil {
		loggers = make(map[*Logger]bool)
	}
	if registries == nil {
		registries = make(map[*Registry]bool)
	}
	// pre order
	registries[root] = true
	for _, l := range root.loggers {
		// avoid dup
		if loggers[l] {
			continue
		}
		loggers[l] = true
		cb(l) // visit
	}
	// dfs
	for _, r := range root.children {
		// avoid cycle
		if registries[r] {
			continue
		}
		walkLogger(r, loggers, registries, cb)
	}
}

func WalkRegistry(root *Registry, cb func(r *Registry)) {
	walkRegistry(root, nil, cb)
}

func walkRegistry(root *Registry, registries map[*Registry]bool, cb func(r *Registry)) {
	// first call
	if registries == nil {
		registries = make(map[*Registry]bool)
	}
	registries[root] = true
	cb(root) // visit
	// dfs
	for _, r := range root.children {
		// avoid dup
		if registries[r] {
			continue
		}
		walkRegistry(r, registries, cb)
	}
}
