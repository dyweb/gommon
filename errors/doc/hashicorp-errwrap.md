# hashicorp/errwrap

- repo: https://github.com/hashicorp/errwrap
- doc: https://godoc.org/github.com/hashicorp/errwrap
- last update: 2014/10/07

Usage

- wrap
- check contains
- get error(s) based on string/type
- walk

````go
func Contains(err error, msg string) bool
func ContainsType(err error, v interface{}) bool
func GetAll(err error, msg string) []error
func GetAllType(err error, v interface{}) []error
func Walk(err error, cb WalkFunc)
func Wrap(outer, inner error) error
func Wrapf(format string, err error) error
type WalkFunc func(error)
type Wrapper interface {
    WrappedErrors() []error
}
````