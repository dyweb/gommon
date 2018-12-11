# Log

- [Design](doc/design)
- [Survey](doc/survey)

## Convention

- library/application MUST have a library/application registry.
- every package MUST have a package level logger.
- logger should only register themselves in registry if it's a long lived, i.e. server

## Usage

See [_examples](_examples)

In your project have a `logutil` package to serve as registry, and import this package in all the other packages
to create package logger so you don't need to manually register them

````go
import (
	"github.com/dyweb/gommon/log"
)

const Project = "github.com/your/project"

var registry = log.NewLibraryRegistry(Project)

func Registry() *log.Registry {
	return &registry
}

func NewPackageLoggerAndRegistry() (*log.Logger, *log.Registry) {
	logger, child := log.NewPackageLoggerAndRegistryWithSkip(Project, 1)
	registry.AddRegistry(child)
	return logger, child
}
````

In package, create package level var called `log` this saves import and also avoid people importing standard log.
Use `dlog` as import alias because `log` is already used.

````go
package server

import (
	dlog "github.com/dyweb/gommon/log"
	"github.com/your/project/logutil"
)

var log, logReg = logutil.NewPackageLoggerAndRegistry()

func foo(file string) {
	// structual way
	log.DebugF("open", dlog.Str("file", file))
	// default handler
	// debug 20180204 open file=test.yml
	// logfmtish handler
	// lvl=debug t=20180204 msg=open file=test.yml
	// json handler
	// {"lvl": "debug", "t": "20180204", "msg": "open", "file": "test.yml"}
	
	// traditional way
	log.Debugf("open %s", file)
	// debug 20180204 open test.yml
	
	// for expensive operation, check before log
	if log.IsDebugEnabled() {
		log.Debug("counter". dlog.Int("counter", CountExpensive()))
	}
}

func bar() {
	// create a new logger based on package level logger
	logger := log.Copy().AddField(dlog.Str("bar", "barbarbar"))
	logger.Info("have extra fields")
}
````