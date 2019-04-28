# Changelog

## v0.0.11

Released on 2019/04/28

[release](https://github.com/dyweb/gommon/releases/tag/0.0.11) [closed issues](https://github.com/dyweb/gommon/issues?q=is%3Aissue+milestone%3A0.0.11+is%3Aclosed)

Breaking change

- [#112](https://github.com/dyweb/gommon/issues/112) Legacy package moved to its own repo [gommon-legacy](https://github.com/dyweb/gommon-legacy)
- [#110](https://github.com/dyweb/gommon/issues/110) Simplify log registry, no longer have tree
- [#97](https://github.com/dyweb/gommon/issues/97) Contains a global registry

New feature

- [#60](https://github.com/dyweb/gommon/issues/60) Use env var `GOMMON_LOG_LEVEL` to set default log level

Removed package

- [#113](https://github.com/dyweb/gommon/issues/113) structure package is moved to legacy

## v0.0.10

[release](https://github.com/dyweb/gommon/releases/tag/0.0.10) [closed issues](https://github.com/dyweb/gommon/issues?q=is%3Aissue+milestone%3A0.0.10+is%3Aclosed)

Breaking change

- [#74](https://github.com/dyweb/gommon/issues/74) Switched to go mod from dep

New feature

- [#76](https://github.com/dyweb/gommon/issues/76) created `errors/errortype` package for inspecting standard library error types, 
and add three common error interfaces `ErrNotFound`, `ErrAlreadyExists`, `ErrNotImplemented`

Fixed

- [#108](https://github.com/dyweb/gommon/issues/108) log package was using `init` to induct package level logger, it is changed to `init.ializers` 

## v0.0.9

[release](https://github.com/dyweb/gommon/releases/tag/0.0.9) [closed issues](https://github.com/dyweb/gommon/issues?q=is%3Aissue+milestone%3A0.0.9+is%3Aclosed)

New feature

- [#66](https://github.com/dyweb/gommon/issues/66) error inspection without checking string, based on [go 2 proposal](https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md)
- [#105](https://github.com/dyweb/gommon/issues/105) `gommon add-build-ignore` used for disable all the legacy packages

Partial

- [#74](https://github.com/dyweb/gommon/issues/74) created `go.mod` so gommon be used by the `replace` directive in other
go mod based projects, gommon itself is still using dep
- [#103](https://github.com/dyweb/gommon/issues/103) style guide, only covered general and is not well formatted

Removed package

- [#68](https://github.com/dyweb/gommon/issues/68) deprecate request package, use httpclient in [go.ice#37](https://github.com/dyweb/go.ice/issues/37)

## v0.0.8

[release](https://github.com/dyweb/gommon/releases/tag/0.0.8) [closed issues](https://github.com/dyweb/gommon/issues?q=is%3Aclosed+milestone%3A0.0.8)

Breaking change in factory functions

- reduce handler interface from 6 methods to 1
- [#88](https://github.com/dyweb/gommon/issues/88) add benchmark
- [#33](https://github.com/dyweb/gommon/issues/33) [#78](https://github.com/dyweb/gommon/issues/33) use log registry to keep tree of logger
  - application logger is now application log registry
  - library logger is now library log registry
- escape string in json handler
- [#87](https://github.com/dyweb/gommon/issues/87) add multi handler for fan out

## v0.0.7

[release](https://github.com/dyweb/gommon/releases/tag/0.0.7)

- deprecated config v1, move it to legacy folder
- move cast into util

## v0.0.6

[release](https://github.com/dyweb/gommon/releases/tag/0.0.6)

- improve doc, fix typo

## v0.0.5

[release](https://github.com/dyweb/gommon/releases/tag/0.0.5)

- remove usage of `pkg/errors`, use `gommon/errors` instead
  - support multiple errors (includes a thread safe version using mutex)
  - `errors.Wrap` only add stack if previous error does not have one

## v0.0.4

[release](https://github.com/dyweb/gommon/releases/tag/v0.0.4)

- add `noodle` package for embed static assets to binary
- allow ignore file and directories using `.noodleignore`

## v0.0.3

[release](https://github.com/dyweb/gommon/releases/tag/v0.0.3)

- add json and cli handler for `gommon/log`
- source line using `runtime`
- add `color` package to util

## v0.0.2 

[release](https://github.com/dyweb/gommon/releases/tag/v0.0.2)

- add generator package and `gommon` cli
  - support generate methods for implementing `log.LoggableStruct` interface
  - use std `flag` package instead of cobra (for now)

