# Application style

Style guide for writing application using gommon or libraries that use gommon

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
