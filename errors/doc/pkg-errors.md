# pkg/errors

- repo: https://github.com/pkg/errors
- doc: http://godoc.org/github.com/pkg/errors
- last update: 2018/01/26

Pending issues 

- [ ] `"Avoid FuncForPC, use CallersFrames` https://github.com/pkg/errors/issues/107
- [ ] use same format verb as go-stack/stack https://github.com/pkg/errors/issues/38
- [ ] Wrap() add withStack only if no cause with stack present https://github.com/pkg/errors/pull/122

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

Internal

- stack call depth is hard coded to 32

````go
type withMessage struct {
	cause error
	msg   string
}

type withStack struct {
	error
	*stack
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	*stack
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
````