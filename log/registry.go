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
	// TODO: it can also be detected using runtime
	Project string
	// Package is detected base on runtime, i.e. github.com/dyweb/gommon/noodle
	Package string
	// Type is specified by user when creating registry
	Type RegistryType
}

func NewPackageRegistryWithSkip(project string, skip int) Registry {
	return Registry{
		identity: newRegistryId(project, PackageRegistry, skip),
	}
}

func NewLibraryRegistry(project string) Registry {
	return Registry{
		identity: newRegistryId(project, LibraryRegistry, 0),
	}
}

func NewApplicationLoggerAndRegistry(project string) (*Logger, *Registry) {
	reg := Registry{
		identity: newRegistryId(project, ApplicationRegistry, 1),
	}
	logger := NewPackageLoggerWithSkip(1)
	reg.AddLogger(logger)
	return logger, &reg
}

func newRegistryId(proj string, tpe RegistryType, skip int) RegistryIdentity {
	// TODO: check if the skip works .... we need another package for testing that
	frame := runtimeutil.GetCallerFrame(2 + skip)
	pkg, _ := runtimeutil.SplitPackageFunc(frame.Function)
	return RegistryIdentity{
		Project: proj,
		Package: pkg,
		Type:    tpe,
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
