# Doc

## Survey

wrap

- [pkg/errors](pkg-errors.md)
- [juju/errgo](juju-errgo.md)
  - support `Mask`, which is not seen in other packages
- [juju/errors](juju-errors.md)
  - same as errgo, and not depending on it, the description is [outdated](https://github.com/juju/errors/pull/10)
- [hashicorp/errwrap](hashicorp-errwrap.md)
  - contains string or type
  - get error by string or type

multi error

- [hashicorp/go-multierror](hashicorp-go-multierror.md)
  - flatten multierror when append
  - ErrorOrNil
    - https://golang.org/doc/faq#nil_error Why is my nil error value not equal to nil?
- [uber/multierr](uber-multierr.md)
  - combine more than one error
  - [x] has a `atomic.Book` for copyNeeded, seems only used for fast path
  
TODO: k8s

- [ ] https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/util/errors/errors.go
- [ ] https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/util/validation/field/errors.go


- implement Formatter interface `Format(f State, c rune)` https://golang.org/pkg/fmt/#Formatter