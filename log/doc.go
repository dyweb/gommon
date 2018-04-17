// Package log provides structured logging with fine grained control
// over libraries using a tree hierarchy of loggers
//
// Conventions
//
// library/application MUST have a library/application logger as registry.
//
// every package MUST have a package level logger.
//
// logger is a registry and can contain children.
//
// instance of struct should have their own logger as children of package logger
package log
