# TODO

## 2018-11-24

- [ ] [#90](https://github.com/dyweb/gommon/pull/90) the long pending benchmark and refactor for adding context (fields to logger instance)

## 2018-05-02

- [x] [#67](https://github.com/dyweb/gommon/issues/67) add example in doc 

## 2018-02-04

- [x] support fields
- [x] support source line, using `runtime` (slow, but requires much less time to implement)

## 2017-12-25

- [x] [#32](https://github.com/dyweb/gommon/issues/32) use runtime.CallersFrame instead of runtime.FuncForPc
- [ ] [#33](https://github.com/dyweb/gommon/issues/33) maintain a tree of logger is harder than I think

## 2017-12-22

- [ ] finish identity, so we can have advanced filtering and print hierarchy
- [ ] allow log file location, it is similar but still different from identity, 
the place you created the logger is not the place you do actual log, 
though they may be very close
- [ ] allow isDebugEnabled (might use RWLock then)

simple usage, one logger in package to rule them all

````go
var log = logutil.NewPkgLogger()

// using package level logger inside function
func foo() {
	log.Info("foo called!")
}

type Auth struct {
}

// use package level logger inside method (in fact still function)
func (a *Auth) check(user string, pwd string) error {
	log.Infof("checking user %s", user)
}
````

fine grained control, also reduce lock contention

- [ ] how to set parent for logger, this may cause cycle import ... (might resolve this in each application/library's logutil)
since all the logger are registered there

````go
var log = logutil.NewPkgLogger()

// var logFoo = logutil.NewFuncLoggerByName(log, "foo")
var logFoo = logutil.NewFuncLogger(log, foo)

func foo() {
	logFoo.Info("I am using my own logger because detail control is needed")
}

type Auth struct {
	l Logger
}

func NewAuth() *Auth {
	return &Auth{
        l: logutil.NewStructLogger(log, Auth),		
	}
}

func (a *Auth) check(user string, pwd string) error {
	a.l.Infof("checking user %s", user)
}
````