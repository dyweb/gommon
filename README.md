# Gommon - Golang common util libraries

[![GoDoc](https://godoc.org/github.com/dyweb/gommon?status.svg)](https://godoc.org/github.com/dyweb/gommon)
[![Build Status](https://travis-ci.org/dyweb/gommon.svg?branch=master)](https://travis-ci.org/dyweb/gommon)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyweb/gommon)](https://goreportcard.com/report/github.com/dyweb/gommon)
[![codebeat badge](https://codebeat.co/badges/8d42a846-f1dc-4a6b-8bd9-5862726ed35d)](https://codebeat.co/projects/github-com-dyweb-gommon-master)

Gommon is a collection of common util libraries that was originally designed for [Ayi](https://github.com/dyweb/Ayi),
and aims to provide a consistent and up to date API for building cli tools and server applications.

It has the following components:

- [Command runner](runner) A processes manager for both short and long running processes.
- [Cast](cast) Convert Golang types
- [Data structure](structure) Bring Set etc. to Golang.
- [Log](log) A log4j style logger, support filtering by package, custom filters etc.
- [Requests](requests) A pythonic wrapper for `net/http`, HTTP for Gopher.

<!--- web server - resource binding (replace go.rice)-->

## Development

- install go
- install [Ayi](https://github.com/dyweb/Ayi)
- install [glide](https://github.com/Masterminds/glide)
- run `Ayi dep-install` or `glide install` to install dependencies.
- run `Ayi test` for test.
- run `Ayi run doc` to view godoc locally, you need to refresh browser and this only works on Linux with `xdg-open`

### Roadmap

- [ ] finish issues transformed from Ayi
- [ ] release 0.0.1

## License

MIT

## Acknowledgement 

Gommon is inspired and in some sense a simplified and unified version of the following awesome libraries

- [logrus](https://github.com/sirupsen/logrus)
- [Cast](https://github.com/spf13/cast)
- [Viper](https://github.com/spf13/viper/)
- [Requests](http://docs.python-requests.org/en/master/)

The name Gommon is suggested by [@arrowrowe](https://github.com/arrowrowe).
The original blog post can be found [here](http://blog.dongyueweb.com/ayi.html).
Thanks all the fellows in [@dyweb](https://github.com/dyweb) for their support.
