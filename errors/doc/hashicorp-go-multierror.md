# hashicorp/go-multierror

- repo: https://github.com/hashicorp/go-multierror
- doc: https://godoc.org/github.com/hashicorp/go-multierror
- last update: 2017/12/04

Usage 

- append
- custom format function
- **flatten**

````go
type Error struct {
	Errors      []error
	ErrorFormat ErrorFormatFunc
}

func (e *Error) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if len(e.Errors) == 0 {
		return nil
	}

	return e
}

func (e *Error) GoString() string {
	return fmt.Sprintf("*%#v", *e)
}

// kind of similar to what is wanted in uber-multierr https://github.com/uber-go/multierr/issues/21
var result *multierror.Error
return result.ErrorOrNil()
````