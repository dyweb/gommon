# Roadmap

## 0.0.8

- [x] test coverage for multiple packages
- [ ] explain internals of some implementation
- [ ] start documenting the style for writing gommon itself, lib using gommon, app using gommon/lib using gommon
- [ ] improve Makefile and dockerized build & test

## 0.0.9

- [ ] more complex error interface, error code
- [ ] organized error types
- [ ] extension for collecting errors using third party services
- [ ] init go mod support, not sure if it will be compatible with dep 

## 0.0.10

- [ ] tree of loggers in use
- [ ] benchmark against other loggers
- [ ] support better logging for errors

## 0.0.11

- [ ] httputil package, merge part of current requests package unix domain sock etc.
- [ ] requests, download and upload file, a curl like example