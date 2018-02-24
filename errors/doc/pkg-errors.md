# pkg/errors

- repo: https://github.com/pkg/errors
- doc: http://godoc.org/github.com/pkg/errors
- last update: 2018/01/26

Pending issues 

- [ ] `"Avoid FuncForPC, use CallersFrames` https://github.com/pkg/errors/issues/107
- [ ] use same format verb as go-stack/stack https://github.com/pkg/errors/issues/38

Usage 

- wrap error
- get cause
- contains stack trace

```go
// wrap errors from standard library
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}

// get cause
type causer interface {
        Cause() error
}
switch err := errors.Cause(err).(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}

type stackTracer interface {
        StackTrace() errors.StackTrace
}
type StackTrace []Frame
```

