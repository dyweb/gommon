# Roadmap

## Up coming

### 0.0.13

Was 0.0.12 but move httpclient from go.ice is needed and improve the error package is going to take a while,
so changed milestone to 0.0.13.

- [ ] align errors with x/errors which will become the default [#109](https://github.com/dyweb/gommon/issues/109)

From 0.0.11

- [ ] requests, download and upload file, a curl like example

From 0.0.10

- [ ] support better logging for errors

From 0.0.9

- [ ] error code
- [ ] organized error types
- [ ] explain internals of some implementation
- [ ] (optional) extension for collecting errors using third party services

## Finished

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
