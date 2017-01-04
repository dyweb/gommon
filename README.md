# Gommon - Golang common util libraries

[![Build Status](https://travis-ci.org/dyweb/gommon.svg?branch=master)](https://travis-ci.org/dyweb/gommon)

<!-- TODO: a better short intro -->
Common utils that are originally from [Ayi](https://github.com/dyweb/Ayi)

- [data structure](structure)
- [log](log)
- [command runner](runner)
- web server
- http client
- resource binding (replace go.rice)

## Data structure

see [data structure](structure)

- Set

## Log

see [log](log)

- filter log by fields, built in support for pkg using `PkgFilter`

## Command runner

see [runner](runner)

- dry run
- run in background (like foreman)
- error handling

## Web server

- static file server (`Ayi web static`) like python's `SimpleHTTPServer`

## Http client

- like python's requests `client.Get('http://api.github.com/').JSON()`
