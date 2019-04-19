# 2019-04-14 Simplify registry

## Issues

- [#110](https://github.com/dyweb/gommon/issues/110) Simplify logger registry

## Background

Although gommon has been keeping a tree of logger, but it is never put into use, there is no UI (commandline/web) that
shows the hierarchy or make use of the hierarchy to generate filtered output, further more forcing the hierarchy makes
the library hard to use. 

Currently registry is like a folder while loggers are like file, a registry has three types, application, library and
package, the original idea is library project can export a registry to allow projects importing this library to control
its behavior directly, however having registry actually make things complex because individual package is independent,
`foo` and `foo/bar` does need to have same hierarchy as their on disk one, further more go does not allow cycle import.

The main problem for logger is fields is still a bit hard to use because unlike other loggers, we only have a `Add`
method to add fields in place and requires using `Copy` to create a copy, it most other loggers they make a copy 
directly when adding fields (some like zap and zerolog serialize the fields when adding them)