# Application style

Style guide for writing application using gommon or libraries that use gommon

## Import

- group import, ref https://github.com/jaegertracing/jaeger/blob/master/CONTRIBUTING.md#imports-grouping
  - std
  - third party
  - packages in project's `lib` folder, they will become third party eventually
  - internal packages
- use `pb` as alias for imported protobuf package
- use `dlog` for `github.com/dyweb/gommon/log` since log is used as package var for package level logger

## Error handling

- when using `log.Fatal`, add a `return` after it, it won't get executed, but it makes the abort logic more obvious
  - if you are calling `panic` then you should not add `return` to avoid `go vet` warning

Good

````go
if cfg.Port < 0 {
    log.Fatalf("invalid port number %d", cfg.Port)
    return
}
server.Serve(cfg.Port)
````

Bad

````go
if cfg.Port < 0 {
    log.Fatalf("invalid port number %d", cfg.Port)
}
server.Serve(cfg.Port)
````
