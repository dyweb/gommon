# v0.0.15 dlog

## TODO

- [ ] list things I want to do to dlog

## Overview

Rename `log` to `dlog` and provider better API for both write and read.

## Motivation

`dyweb/log` spent a lot of time on write performance instead of usability.
It has a structured logging API, but I find it hard to use for both cli and server apps.

## Implementation

- [ ] rename `log` package to `dlog`
- [ ] add the ability to read the log it generates
- [ ] reduce size of interface
  - [ ] remove `PanicX` simply call `ErrorX` and `panic` 
