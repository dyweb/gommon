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

Go logging library rarely keep a tree hierarchy, it is not the case in Java.
In [Solr](../survey/solr.md) where log4j1 is used, relationship is kept inside loggers,
in [log4j 2](../survey/log4j.md) they introduced logger config with hierarchy, it is determined by logger name and dot.
(TODO: is this hierarchy applied every time when log or only once on configuration change)