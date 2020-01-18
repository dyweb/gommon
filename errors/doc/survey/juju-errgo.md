# juju/errgo

- repo: https://github.com/juju/errgo
- doc: https://godoc.org/github.com/juju/errgo
- last update: 2014/09/25

Usage

- [ ] Mask ... https://github.com/juju/errgo#func--mask
- cause

````go
func Any(error) bool
func Cause(err error) error
// return stack etc. in a single string
func Details(err error) string
func Mask(underlying error, pass ...func(error) bool) error
func WithCausef(underlying, cause error, f string, a ...interface{}) error
````