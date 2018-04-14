# Gommon - Go common libraries

[![GoDoc](https://godoc.org/github.com/dyweb/gommon?status.svg)](https://godoc.org/github.com/dyweb/gommon)
[![Build Status](https://travis-ci.org/dyweb/gommon.svg?branch=master)](https://travis-ci.org/dyweb/gommon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyweb/gommon)](https://goreportcard.com/report/github.com/dyweb/gommon)
[![codebeat badge](https://codebeat.co/badges/8d42a846-f1dc-4a6b-8bd9-5862726ed35d)](https://codebeat.co/projects/github-com-dyweb-gommon-master)
[![Sourcegraph](https://sourcegraph.com/github.com/dyweb/gommon/-/badge.svg)](https://sourcegraph.com/github.com/dyweb/gommon?badge)

Gommon is a collection of common util libraries written in Go.

It has the following components:

- [config](config) A YAML config reader with template support
- [errors](errors) Wrap error and multi error
- [generator](generator) Render go template, generate methods for logger interface based on `gommon.yml`
- [log](log) A Javaish logger for Go, application can control library and set level for different pkg via config or flag
- [noodle](noodle) Embed static assets for web application with `.noodleignore` support
- [requests](requests) A pythonic wrapper for `net/http`, HTTP for Gopher.
- [cast](cast) Convert Golang types
- [structure](structure) Bring data structure like Set etc. to Golang.

Legacy

- [log v1](legacy/log) A logrus like structured logger
- [Runner](legacy/runner) A os/exec wrapper


## Dependencies

Currently we only have one non standard library dependencies, see [Gopkg.lock](Gopkg.lock)

- [go-yaml/yaml](https://github.com/go-yaml/yaml) for read config written in YAML
  - we don't need most feature of YAML, and want to have access to the parser directly to report which line has incorrect semantic (after checking it in application).
    - might write one in [ANTLR](https://github.com/antlr/antlr4)
  - we also have a DSL work in progress [RCL: Reika Configuration Language](https://github.com/at15/reika/issues/49), which is like [HCL](https://github.com/hashicorp/hcl2)

Removed 

- [pkg/errors](https://github.com/pkg/errors) for including context in error
  - removed in [#59](https://github.com/dyweb/gommon/pull/59)
  - replaced by `gommon/errors`

<!-- no, we are using the standard flag package ... -->
<!-- For command line util, we are using [spf13/cobra](https://github.com/spf13/cobra), it is more flexible than [ufrave/cli](https://github.com/urfave/cli) -->

## Development

- install go https://golang.org/
- install dep https://github.com/golang/dep
- `dep ensure`
- `make test`

### Roadmap

0.0.6 

- [ ] better documentation, add examples
- [ ] explain internals of some implementation

0.0.7

- [ ] organized error types 
- [ ] extension for collecting errors using third party services

0.0.8

- [ ] tree of loggers in use
- [ ] benchmark against other loggers
- [ ] support better logging for errors

## License

MIT

## Contribution

Currently, gommon is in a very violate state, please open issues after it becomes stable 

## Acknowledgement & Improvement

Gommon is inspired by the following awesome libraries, 
most gommon packages have much less (not more) features and a few improvements 
compared to packages it modeled after.

log

- [sirupsen/logrus](https://github.com/sirupsen/logrus) for structured logging 
  - log v1 is entirely modeled after logrus, entry contains log information with methods like `Info`, `Infof`
- [apex/log](https://github.com/apex/log) for log handlers
  - log v2's handler is inspired by apex/log, but we didn't use entry and chose to pass multiple parameters to explicitly state what a handler should handle
- [uber-go/zap](https://github.com/uber-go/zap) for serialize log fields without using `fmt.Sprintf` and use `strconv` directly
  - we didn't go that extreme as Zap or ZeroLog for zero allocation, performance is not measured

config

- [spf13/cast](https://github.com/spf13/cast) for cast, it is used by Viper
- [spf13/viper](https://github.com/spf13/viper/) for config
  - looking up config via string key makes type system useless, so we always marshal entire config file to a single struct
    - it also makes refactor easier
  - we only use YAML, might add RCL, won't support TOML

requests

- [Requests](http://docs.python-requests.org/en/master/) for requests
- [hashicorp/go-cleanhttp](https://github.com/hashicorp/go-cleanhttp) for using non default http transport and client

generator

- [benbjohnson/tmpl](https://github.com/benbjohnson/tmpl) for go template generator
  - first saw it in [influxdata/influxdb](https://github.com/influxdata/influxdb/blob/master/tsdb/engine/tsm1/encoding.gen.go.tmpl)
  - we put template data in `gommon.yml`, so we don't need to pass data as json via cli

noodle

- [GeertJohan/go.rice](https://github.com/GeertJohan/go.rice)
  - we implemented `.gitignore` like [feature](https://github.com/at15/go.rice/issues/1) but the upstream didn't respond for the [feature request #83](https://github.com/GeertJohan/go.rice/issues/83)
  - we put data into generated code file while go.rice and append zip to existing go binary
  
errors

- [pkg/errors](https://github.com/pkg/errors) it can not introduce breaking change, but `WithMessage` and `WithStack` is annoying
  - see [#54](https://github.com/dyweb/gommon/issues/54) and [errors/doc](errors/doc) about other error packages
  - https://github.com/pkg/errors/pull/122 for check existing stack before attach new one
- [uber-go/multierr#21]( https://github.com/uber-go/multierr/issues/21) for return bool after append
- [hashicorp/go-multierror](https://github.com/hashicorp/go-multierror) for `ErrorOrNil`

## About

It was part of [Ayi](https://github.com/dyweb/Ayi) and split out for wider use.
The name Gommon is suggested by [@arrowrowe](https://github.com/arrowrowe).
The original blog post can be found [here](http://blog.dongyueweb.com/ayi.html).
Thanks all the fellows in [@dyweb](https://github.com/dyweb) especially [@gaocegege](https://github.com/gaocegege) for their support in early development.
