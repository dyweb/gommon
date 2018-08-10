# Errors

A error package supports wrapping and multi error

## Issues

- init https://github.com/dyweb/gommon/issues/54

## Implementation

- `Wrap` checks if this is already a `WrappedErr`, if not, it attach stack

## Survey

- Wrap error with context
  - [pkg/errors](doc/pkg-errors.md)
  - [juju/errors](doc/juju-errors.md)
  - [hashicorp/errwrap](doc/hashicorp-errwrap.md)
- Multi errors
  - [hashicorp/go-multierror](doc/hashicorp-go-multierror.md)
  - [uber/multierr](doc/uber-multierr.md)