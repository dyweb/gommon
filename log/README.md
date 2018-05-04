# Log

## Convention

- library/application MUST have a library/application logger as registry.
- every package MUST have a package level logger.
- logger is a registry and can contain children.
- instance of struct should have their own logger as children of package logger

## Usage

````go
package logutil

import (
	"github.com/dyweb/gommon/log"
)

var Registry = log.NewLibraryLogger()

func NewPackageLogger() *log.Logger {
	l := log.NewPackageLoggerWithSkip(1)
	Registry.AddChild(l)
	return l
}
````


````go
var log = logutil.NewPackageLogger()

func foo() {
	// structual way
	log.DebugF("open", dlog.Fields{"file": file})
	// default handler
	// debug 20180204 open file=test.yml
	// logfmtish handler
	// lvl=debug t=20180204 msg=open file=test.yml
	// json handler
	// {"lvl": "debug", "t": "20180204", "msg": "open", "file": "test.yml"}
	// traditional way
	log.Debugf("open %s", file)
	// debug 20180204 open test.yml
	// a mixed way, this would lose hint from IDE for printf placeholders
	log.DebugFf(dlog.Fields{"file": file}, "open with error %v", err)
	// default handler
	
	// for expensive operation, check before log
	if log.IsDebugEnabled() {
		log.Debug("counter". dlog.Fields{"counter": CountExpensive()})
	}
}
````