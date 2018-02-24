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