# Errors

A error package supports wrapping and multi error

## Issues

- init https://github.com/dyweb/gommon/issues/54

## Implementation

- `Wrap` checks if this is already a `WrappedErr`, if not, it attach stack

## Survey

- Wrap error with context (stack, cause etc.)
  - [pkg/errors](doc/pkg-errors.md)
    - Wrap() add withStack only if no cause with stack present https://github.com/pkg/errors/pull/122/
  - [juju/errors](doc/juju-errors.md)
  - [hashicorp/errwrap](doc/hashicorp-errwrap.md)
  - [ ] https://godoc.org/github.com/gorilla/securecookie#Error
- Multi errors
  - [hashicorp/go-multierror](doc/hashicorp-go-multierror.md)
  - [uber/multierr](doc/uber-multierr.md)
  - [ ] https://godoc.org/github.com/gorilla/securecookie#MultiError