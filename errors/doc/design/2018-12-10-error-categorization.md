# 2018-12-10 Error categorization

## Issues

- [#62](https://github.com/dyweb/gommon/issues/62) Support fmt.Formatter interface
- [#73](https://github.com/dyweb/gommon/issues/73) Human readable suggestions for possible solutions
- [#76](https://github.com/dyweb/gommon/issues/76) A more complex error interface

## Background

Originally I was using [pkg/errors](https://github.com/pkg/errors), it works fine.
I wrap all my errors and many library I depend on also use it, there is not import issues.
However, my handling logic remain the same, have a `log.Fatal(err)` if it's a cli
or `w.Write([]byte(err.Error()))` if it's an http server. 
I don't generate helpful suggestion like `you want config.yaml but only config.yml is found in current folder`.
for my library and app users (which sadly just myself most of the time).

The are two reasons for that, First I don't know what errors will the standard library return,
there are detailed godoc, but I want them in one page so I don't need to jump around the doc.
Second I am just wrapping for the sake of wrapping, add things like `error read file` helps
human a bit when it got logged at last, but it's not a context that can be used by code, 
in the end all the read file error are treated the same regardless why I read the file.

However, the direct reason for me to stop using pkg/errors are there are pending issues and 
I want to have multi error. pkg/errors is not going to add more complexity,
i.e. [#75](https://github.com/pkg/errors/issues/75) if you wrap a wrapped error, 
it is going to have duplicated stack with less trace, the more you wrap, the less you have after first Cause (UnWrap).
The PR [#122](https://github.com/pkg/errors/pull/122) to fix that has been there for almost 2 years.
As for multi error, [hashicorp](https://github.com/hashicorp/go-multierror) and [uber-go/multierr](https://github.com/uber-go/multierr)
are good example, they also have pending issues like return bool for Append [#21](https://github.com/uber-go/multierr/issues/21)

Thus I reinvent the wheel and gommon/errors is a mix of all of them, it contains error wrapping and multi error,
and implemented some of their pending issues like duplicated wrap and return append result in multi error.

<!-- the above section is quite messy ... might need to rewrite it ... -->

### Current implementation

The implementation is pretty simple, the most complex part, print call stack, is not even implemented [#62](https://github.com/dyweb/gommon/issues/62),
it does not support `%+v` formatter. 

Wrap error is just an interface with one implementation which store a message as context and stack trace

````go
type WrappedError struct {
	msg   string
	cause error
	stack *Stack
}
````

Multi error has two implementation, one is thread safe by using a mutex, one is just a slice

````go
// multiErr is NOT thread (goroutine) safe
type multiErr struct {
	errs []error
}

// multiErrSafe is thread safe
type multiErrSafe struct {
	mu   sync.Mutex
	errs []error
}
````

## Goals

There are many things I planned to do, some are too heavy that I think a `errorx` library is needed,
will have a `logx` package as well

- categorized errors
  - standard library errors
  - dependency errors (not my application/library)
  - network, permanent and temp
  - client/server error, client should know it's a server error or client itself figured out something wrong before/after calling server
  - encoding/decoding
  - invalid input i.e. wrong file name
  - implementation, assert/sanity check failed at runtime and you don't want to panic
- allow traverse error chain in code
- try to align w/ [go2 proposals](https://go.googlesource.com/proposal/+/master/design/go2draft.md)
- error collection and aggregation for real-time and future analysis

## Plan

First is to read [go2draft](https://go.googlesource.com/proposal/+/master/design/go2draft.md) again, 
I have already forgot most of them, luckily I did take [notes](https://github.com/at15/papers-i-read/commit/a44b13fd8e1bcc51135ebe128f391ebc674904c4)
The four ways of error handling and the problem user want to ask when there is an RPC error 
is really what I always wanted to say but didn't figure out how to (sigh).

Second go over all the refers I put in [#76](https://github.com/dyweb/gommon/issues/76) and [go.ice#12](https://github.com/dyweb/go.ice/issues/12)

- TiDB https://github.com/pingcap/tidb/tree/master/terror
- gRPC https://github.com/grpc-ecosystem/grpc-opentracing/blob/master/go/otgrpc/errors.go has error class
- docker https://github.com/moby/moby/blob/master/errdefs/doc.go `Errors that cross the package boundary should implement one (and only one) of these interfaces.`
- teleport (ssh servers) https://github.com/gravitational/trace/blob/master/errors.go concrete structs (different from docker's errdefs)
- dgraph https://github.com/dgraph-io/dgraph/blob/master/x/error.go the `func Check2(_ interface{}, err error)` is epic, I want `func Check3(_ interface{}, _ interface{}, err error)`
- [ ] error collection like sentry, it's good to ship all your errors like logs into a central place for cross application/node analysis
(if you can afford it)

Third lines up the following

- list all the common errors from standard library, i.e. file, http, encoding
  - common error message (amazing error messages and where to find them)
- what categories will be included in `errors`
- ~~three~~ some real world use cases
  - ayi, a command line productivity tool that read config, make some API call, run other processes
  - libtsdb (well the project is dead but ...) a library that do encoding and calling servers
  - benchhub, a web service that talks with many other error services
- what's in `errors` what's in `errorx`