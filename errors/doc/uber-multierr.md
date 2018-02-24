# uber-go/multierr

- repo: https://github.com/uber-go/multierr
- doc: https://godoc.org/go.uber.org/multierr

Pending

- [ ] zap integration https://github.com/uber-go/multierr/issues/23
- [ ] return two value for append https://github.com/uber-go/multierr/issues/21

TODO

- [ ] it seems it can't handle concurrent access? it has an atomic value `copyNeeded`

Usage

- Combine and Append

````go
multierr.Combine(
	reader.Close(),
	writer.Close(),
	conn.Close(),
)
err = multierr.Append(reader.Close(), writer.Close())
````

````go
type multiError struct {
	copyNeeded atomic.Bool
	errors     []error
}

````

- append has fast path, which use `atomic.Bool` to determine is copy is needed
- flatten multierrors

````go
func Append(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}

	if _, ok := right.(*multiError); !ok {
		if l, ok := left.(*multiError); ok && !l.copyNeeded.Swap(true) {
			// Common case where the error on the left is constantly being
			// appended to.
			errs := append(l.errors, right)
			return &multiError{errors: errs}
		} else if !ok {
			// Both errors are single errors.
			return &multiError{errors: []error{left, right}}
		}
	}

	// Either right or both, left and right, are multiErrors. Rely on usual
	// expensive logic.
	errors := [2]error{left, right}
	return fromSlice(errors[0:])
}
````