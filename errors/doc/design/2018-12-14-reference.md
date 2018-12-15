# 2018-12-14 Reference

This doc go over the references mentioned in [error-categorization](2018-12-10-error-categorization.md)

## Go 2

- notes: https://github.com/at15/papers-i-read/commit/a44b13fd8e1bcc51135ebe128f391ebc674904c4
- source: https://go.googlesource.com/proposal/+/master/design/go2draft-error-values-overview.md
 
TODO

- [ ] it mentioned https://github.com/spacemonkeygo/errors which has class etc.
  - it allows attach key value pairs to error
  - defines common set of error https://github.com/spacemonkeygo/errors/blob/master/errors.go#L563-L600
- [ ] https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html by rob pike
  - defines error kinds https://godoc.org/upspin.io/errors#Kind
  - https://godoc.org/upspin.io/errors#MarshalError to transfer error across the wire

https://go.googlesource.com/proposal/+/master/design/go2draft-error-printing.md

- error printing is only for read by human
- it print trace of error (w/ or w/o stack?)
  - [ ] now in gommon/errors we only do wrapping in first error, however the call stack of init error should be different from the wrapping stack
- mentioned list of error (which is multi error)

https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md

- `Is` is just same as upspin, and I think it only works for sentinel errors like `io.EOF`

````go
func Is(err, target error) bool {
	for {
		if err == target {
			return true
		}
		wrapper, ok := err.(Wrapper)
		if !ok {
			return false
		}
		err = wrapper.Unwrap()
		if err == nil {
			return false
		}
	}
}
````

- `As` requires contracts ... (pass type as a parameter)
  - without polymorphism
  - [ ] I didn't quite get this part ..
 
````go
func As(type E)(err error) (e E, ok bool) {
	for {
		if e, ok := err.(E); ok {
			return e, true
		}
		wrapper, ok := err.(Wrapper)
		if !ok {
			return e, false
		}
		err = wrapper.Unwrap()
		if err == nil {
			return e, false
		}
	}
}
````

````go
// instead of pe, ok := err.(*os.PathError)
var pe *os.PathError
if errors.AsValue(&pe, err) { ... pe.Path ... }
````

## Upspin

https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

- didn't use stack trace for error, show operational trace

> There is a tension between making errors helpful and concise for the end user versus making them expansive and analytic for the implementer. 
Too often the implementer wins and the errors are overly verbose, to the point of including stack traces or other overwhelming detail

> Upspin's errors are an attempt to serve both the users and the implementers. 
The reported errors are reasonably concise, concentrating on information the user should find helpful. 
But they also contain internal details such as method names an implementer might find diagnostic but not in a way that overwhelms the user. 
In practice we find that the tradeoff has worked well

> In contrast, a stack trace-like error is worse in both respects. 
The user does not have the context to understand the stack trace, 
and an implementer shown a stack trace is denied the information that could be presented 
if the server-side error was passed to the client. This is why Upspin error nesting behaves as an operational trace, 
showing the path through the elements of the system, rather than as an execution trace, showing the path through the code. 
The distinction is vital

- `func Is(kind Kind, err error) bool` I think is almost same as the go 2 proposal `Is`, compare through the cause chain
- `func Match(template, err error) bool` is almost same as go 2 proposal `func As(type E)(err error) (e E, ok bool) {}`

````go
// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

// Match only operates on upspin's Error type, so it does not need the polymorphism like go 2
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	// un wrap and compare etc.
}
````