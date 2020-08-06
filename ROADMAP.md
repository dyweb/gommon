# Roadmap

NOTE: it is being moved to [milestones](doc/milestones)

## Up coming

### 0.0.14

- [ ] dcli
- [ ] wait package, similar to the polling package in k8s
  - [ ] have retry as alias and provides backoff
  - [ ] allow use wait for container

From 0.0.13

- [ ] log UI, parse log output etc.

From 0.0.11

- [ ] requests, download and upload file, a curl like example

From 0.0.10

- [ ] support better logging for errors

From 0.0.9

- [ ] error code
- [ ] organized error types
- [ ] explain internals of some implementation
- [ ] (optional) extension for collecting errors using third party services

- [x] support `.gommonignore`, used by `gommon generate` on `_legacy` folder
- [x] clean up `.ignore` file support, move it to `fsutil` package
- [x] clean up `go.mod` by using a separated `go.mod` for `gommon` binary, thus removing cobra, viper, etcd.
- [x] min/max for integer https://github.com/dyweb/gommon/issues/123
- [x] create test container https://github.com/dyweb/gommon/issues/124

## Finished

### 0.0.13

- [x] align errors with x/errors which will become the default [#109](https://github.com/dyweb/gommon/issues/109)

### 0.0.12

- [x] move httpclient from go.ice

### 0.0.11

- [x] simplify log package [#110](https://github.com/dyweb/gommon/issues/110)
- [x] move deprecated package to its own repo [#112](https://github.com/dyweb/gommon/issues/112)

### 0.0.10

- [x] switch from dep to go mod

### 0.0.9

- [x] more complex error interface
- [x] start documenting the style for writing gommon itself, lib using gommon, app using gommon/lib using gommon
- [x] improve Makefile and dockerized build & test
- [x] init go mod
- [x] httputil package, merge part of current requests package unix domain sock etc.

### 0.0.8

- [x] test coverage for multiple packages
- [x] tree of loggers in use
- [x] benchmark against other loggers
