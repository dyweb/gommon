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

- I was thinking it is using `strings.Contains` but it is actually comparing message directly ...

````go
// GetAll gets all the errors that might be wrapped in err with the
// given message. The order of the errors is such that the outermost
// matching error (the most recent wrap) is index zero, and so on.
func GetAll(err error, msg string) []error {
	var result []error

	Walk(err, func(err error) {
		if err.Error() == msg {
			result = append(result, err)
		}
	})

	return result
}
````

````go
// GetAllType gets all the errors that are the same type as v.
//
// The order of the return value is the same as described in GetAll.
func GetAllType(err error, v interface{}) []error {
	var result []error

	var search string
	if v != nil {
		search = reflect.TypeOf(v).String()
	}
	Walk(err, func(err error) {
		var needle string
		if err != nil {
			needle = reflect.TypeOf(err).String()
		}

		if needle == search {
			result = append(result, err)
		}
	})

	return result
}
````