# Gommon - Go common libraries

<h1 align="center">
	<br>
	<img width="100%" src="doc/media/gommon.png" alt="gommon">
	<br>
	<br>
	<br>
</h1>

[![GoDoc](https://godoc.org/github.com/dyweb/gommon?status.svg)](https://godoc.org/github.com/dyweb/gommon)
[![Build Status](https://travis-ci.org/dyweb/gommon.svg?branch=master)](https://travis-ci.org/dyweb/gommon)
[![codecov](https://codecov.io/gh/dyweb/gommon/branch/master/graph/badge.svg)](https://codecov.io/gh/dyweb/gommon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyweb/gommon)](https://goreportcard.com/report/github.com/dyweb/gommon)
[![loc](https://tokei.rs/b1/github/dyweb/gommon)](https://github.com/dyweb/gommon)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdyweb%2Fgommon.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdyweb%2Fgommon?ref=badge_shield)

Gommon is a collection of common util libraries written in Go.

- [errors](errors) error wrapping, inspection, multi error (error list), common error types
- [log](log) per package logger with [reasonable performance](log/_benchmarks/README.md)
- [noodle](noodle) embed static assets for web application with `.noodleignore` support
- [generator](generator) render go template, generate methods for logger interface based on `gommon.yml`
- [util](util) wrappers for standard libraries

It has little third party dependencies, only [go-yaml/yaml](https://github.com/go-yaml/yaml) in [util/cast](util/cast),
[go-shellquote](github.com/kballard/go-shellquote) in [generator](generator),
other dependencies like cobra are only for cli, see [go.mod](go.mod).

## Development

- requires go1.12+. go1.11.x should work as well, the Makefile set `GO111MODULE=on` so you can use in GOPATH
- `make help`
- [Directory layout](directory.md)

## License

MIT

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdyweb%2Fgommon.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdyweb%2Fgommon?ref=badge_large)

## Contribution

Currently, gommon is in a very volatile state, please open issues after it becomes stable.

## About

Gommon is inspired by many existing libraries, attribution and comparision can be found in [doc/attribution](doc/attribution.md).

Gommon was part of [Ayi](https://github.com/dyweb/Ayi) and split out for wider use.
The name Gommon is suggested by [@arrowrowe](https://github.com/arrowrowe).
The original blog post can be found in dongyue web's [blog](http://blog.dongyueweb.com/ayi.html).
Thanks all the folks in [@dyweb](https://github.com/dyweb)
especially [@gaocegege](https://github.com/gaocegege) for their support in early development.