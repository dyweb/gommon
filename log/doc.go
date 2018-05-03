/*
Package log provides structured logging with fine grained control
over libraries using a tree hierarchy of loggers

Conventions

1. no direct use of the log package, MUST create new logger.

2. library/application MUST have a library/application logger as their registry.

3. every package MUST have a package level logger as child of the registry, normally defined in pkg.go

4. logger is a registry and can contain children.

5. instance of struct should have their own logger as children of package logger
*/
package log
