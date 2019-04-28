# 2019-04-14 Simplify registry

## Issues

[All issues in 0.0.11](https://github.com/dyweb/gommon/issues?q=is%3Aopen+is%3Aissue+project%3Adyweb%2Fgommon%2F2+milestone%3A0.0.11)

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

Go logging library rarely keep a tree hierarchy, normally people use a global logger, i.e. `log`, `logrus`, `glog`
It is not the case in Java, in [Solr](../survey/solr.md) where log4j1 is used, relationship is kept inside loggers,
in [log4j 2](../survey/log4j.md) they introduced logger config with hierarchy, it is determined by logger name and dot.
(TODO: is this hierarchy applied every time when log or only once on configuration change)

A main problem of having hierarchy is config becomes harder, things like using debug level becomes changing level to debug
for all loggers based on registry, for those didn't register in registry, they will never get updated
([#97](https://github.com/dyweb/gommon/issues/97) register to global logger registry by default)

## Design

Why we need to keep track of logger

- allow individual logger to have and change to different behaviors (i.e. handler) at runtime, otherwise we can just point to a global config in each logger or simply have a global logger
- when update logger config, it should be propagate to all loggers
  - logger config is read heavy, so we pay cost when update, there is very few config update, but config read happens every time when log

What kind of logger should get tracked

- long lived (span the entire application lifecycle), i.e. package level logger
- very few instances, i.e. singleton/limited struct loggers
- has value in customization (i.e. change level)

What kind of logger should NOT get tracked

- created as copy for adding context for current func, i.e. a method logger from struct logger, `st.logger.WithMethod()`
- too many of them, will blow up registry (unless deregister is provided, but this requires remove one element from slice)

Registry

- only one type of registry, package registry, no more application & library, main package is application,
your meta package (import and alias all your subpackages) is library (saw in controller runtime [alias.go](https://github.com/kubernetes-sigs/controller-runtime/blob/master/alias.go))

## Implementation

````go
// registry
var globalRegistryGroup

type registryGroup {
	mu sync.Mutex
	registris map[string]*Registry // key is package name
}

func NewRegistry() {
	reg := Registry{id: caller(), default: xxx, loggers:[]{xxx}}
    globalRegistryGroup.add(id.Pkg, reg)
	return reg
}

func WalkRegistry(fn func(r *Registry)) {
	// can first sort the keys to have a consistent walk order
	// ...
}

func WalkLogger() {
	// walk into logger of individual registry
}

// usage
// in pkg.go
var logReg = dlog.NewRegistry()
var log = logReg.Default()

// in struct impl
func NewController() *Controller{
	c := &Controller{}
	c.logger = logReg.NewLogger(?id)
	return c
}
````