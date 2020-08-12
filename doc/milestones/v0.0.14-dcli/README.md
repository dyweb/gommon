# Gommon v0.0.14 dcli

## TODO

- [ ] merge w/ existing design doc in [dcli/doc/design](../../../dcli/doc/design)
- [ ] split up features
- [ ] list implementation order
 
## Overview

A commandline application builder `dcli` that replaces [spf13/cobra](https://github.com/spf13/cobra).
Minor fix to update small util packages e.g. `mathutil`, `stringutil`, `envutil`.

## Motivation

`dcli`

- less dependencies
- more customization
- type safe
- learn from existing command line builders in different languages, e.g. clap (w/ structopt)

## Implementation

- [x] define `Command` interface
- [x] define a common `Command` implementation so user don't need to create new struct most of the time
  - [x] `Cmd` is the struct
- [x] look up sub command `binary foo bar boar`
  - assume there is no flag, e.g. no `binary foo --bla bar`
  - assume there is no position argument, e.g. no `binary foo arg1 arg2`
- [ ] show help message
  - [x] define an error for command whose sole purpose is showing help, e.g. `git`
  - [ ] a global error handler (or per command)
  - [ ] show formatted error message
  - [ ] dump error message as HTML (there is no flag support, how to do that now, env?)
- [ ] rename `gommon` binary package to `gom` or whatever way to make it's easier to install the binary w/ `go get`
- [ ] remove `cobra` from `gommon` binary dependency

## Specs

- support git style flag and subcommand
- use `dcli` for `gommon` command

## Features

### Sub command

Description

Support subcommand like `git clone`, and show help message for available sub command.
For simplicity in this feature we don't consider flag and position argument.

Components

- `help`
  - list sub command in help messages, and provide short description
  - print a html page if there are too many things for terminal and user prefer clicking around
  - print a text file so user can open it in vim and search it
  
### Flag

Description

Flag can show up in multiple places

- after binary name `foo --verbose=true`
- after sub command `foo bar --target==127.0.0.1:3530`

We should follow spf13/cobra where you can define a flag in parent command and shared by all child commands.
i.e. `PersistentFlag`