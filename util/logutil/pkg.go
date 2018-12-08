// Package logutil is a registry of loggers, it is required for all lib and app that use gommon/log.
// You should add the registry as child of your library/application's child if you want to control gommon libraries
// logging behavior
package logutil

import (
	"github.com/dyweb/gommon/log"
)

const Project = "github.com/dyweb/gommon"

var registry = log.NewLibraryRegistry(Project)

func Registry() *log.Registry {
	return &registry
}

func NewPackageLoggerAndRegistry() (*log.Logger, *log.Registry) {
	child := log.NewPackageRegistryWithSkip(Project, 1)
	logger := log.NewPackageLoggerWithSkip(1)
	child.AddLogger(logger)
	registry.AddRegistry(&child)
	return logger, &child
}
