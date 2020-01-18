# Survey of go error handling libraries

- Wrap error with context (stack, cause etc.)
  - [pkg/errors](pkg-errors.md)
    - Wrap() add withStack only if no cause with stack present https://github.com/pkg/errors/pull/122/
  - [juju/errors](juju-errors.md)
  - [hashicorp/errwrap](hashicorp-errwrap.md)
  - [ ] https://godoc.org/github.com/gorilla/securecookie#Error
- Multi errors
  - [hashicorp/go-multierror](hashicorp-go-multierror.md)
  - [uber/multierr](uber-multierr.md)
    - Wish we could tell if the append happened  https://github.com/uber-go/multierr/issues/21
  - [ ] https://godoc.org/github.com/gorilla/securecookie#MultiError