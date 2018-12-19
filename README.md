# Gommon - Go common libraries

[![GoDoc](https://godoc.org/github.com/dyweb/gommon?status.svg)](https://godoc.org/github.com/dyweb/gommon)
[![Build Status](https://travis-ci.org/dyweb/gommon.svg?branch=master)](https://travis-ci.org/dyweb/gommon)
[![codecov](https://codecov.io/gh/dyweb/gommon/branch/master/graph/badge.svg)](https://codecov.io/gh/dyweb/gommon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyweb/gommon)](https://goreportcard.com/report/github.com/dyweb/gommon)
[![codebeat badge](https://codebeat.co/badges/8d42a846-f1dc-4a6b-8bd9-5862726ed35d)](https://codebeat.co/projects/github-com-dyweb-gommon-master)
[![Sourcegraph](https://sourcegraph.com/github.com/dyweb/gommon/-/badge.svg)](https://sourcegraph.com/github.com/dyweb/gommon?badge)
[![](https://tokei.rs/b1/github/dyweb/gommon)](https://github.com/dyweb/gommon)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdyweb%2Fgommon.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdyweb%2Fgommon?ref=badge_shield)

Gommon is a collection of common util libraries written in Go.

It has the following components:

- [errors](errors) error wrapping, inspection, multi error (error list), common error types
- [log](log) fine grained level control and reasonable performance
- [noodle](noodle) embed static assets for web application with `.noodleignore` support
- [generator](generator) render go template, generate methods for logger interface based on `gommon.yml`
- [structure](structure) data structure like Set etc. to go
- [util](util) small utils over standard libraries utils

Deprecating

- [requests](requests) A pythonic wrapper for `net/http`, HTTP for Gopher

Legacy

- [config v1](config) A YAML config reader with template support
- [log v1](legacy/log) A logrus like structured logger
- [Runner](legacy/runner) A os/exec wrapper

## Dependencies

Currently we only have one non standard library dependencies (cmd and examples are not considered), see [Gopkg.lock](Gopkg.lock)

- [go-yaml/yaml](https://github.com/go-yaml/yaml) for read config written in YAML
  - we don't need most feature of YAML, and want to have access to the parser directly to report which line has incorrect semantic (after checking it in application).
    - might write one in [ANTLR](https://github.com/antlr/antlr4)
  - we also have a DSL work in progress [RCL: Reika Configuration Language](https://github.com/at15/reika/issues/49), which is like [HCL](https://github.com/hashicorp/hcl2)

Removed

- [pkg/errors](https://github.com/pkg/errors) for including context in error
  - removed in [#59](https://github.com/dyweb/gommon/pull/59)
  - replaced by `gommon/errors`

## Development

- install go https://golang.org/
- install dep https://github.com/golang/dep
- `dep ensure`
- `make test`

## License

MIT

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdyweb%2Fgommon.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdyweb%2Fgommon?ref=badge_large)

## Contribution

Currently, gommon is in a very violate state, please open issues after it becomes stable.

## About

Gommon is inspired by many existing libraries, attribution and comparision can be found in [doc/attribution](doc/attribution.md).

Gommon was part of [Ayi](https://github.com/dyweb/Ayi) and split out for wider use.
The name Gommon is suggested by [@arrowrowe](https://github.com/arrowrowe).
The original blog post can be found in dongyue web's [blog](http://blog.dongyueweb.com/ayi.html).
Thanks all the folks in [@dyweb](https://github.com/dyweb)
especially [@gaocegege](https://github.com/gaocegege) for their support in early development.