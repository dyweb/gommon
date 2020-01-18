# Errors

[![GoDoc](https://godoc.org/github.com/dyweb/gommon/errors?status.svg)](https://godoc.org/github.com/dyweb/gommon/errors)

A replacement for `errors` and `pkg/errors`. Supports error wrapping, inspection and multi errors.

- [Survey](doc/survey)

## Issues

- [#54](https://github.com/dyweb/gommon/issues/54) init 

## Implementation

- [doc/design](doc/design) contains some survey and design choices we made
- `Wrap` checks if this is already a `WrappedErr`, if not, it attaches stack
- `MultiErr` keep a slice of errors, the thread safe version use mutex and returns copy of slice when `Errors` is called

## References and Alternatives

- [golang/xerrors](https://github.com/golang/xerrors) Official error wrapping `fmt.Errorf("oh my %w", err)`
- [rotisserie/eris](https://github.com/rotisserie/eris)
- [go-rewrap-errors](https://github.com/xdg-go/go-rewrap-errors) Conver pkg/errors to standard library wrapping