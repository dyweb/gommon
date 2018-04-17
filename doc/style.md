# Coding Style

The coding style can be split into two parts, application and library,
gommon is mainly a set of libraries, but it also contains a command line application with same name `gommon`.

## Folder structure

see [directory](directory.md)

## Documentation

MUST cover the following

- [ ] convention, i.e. variable names, error handling etc.
- [ ] internal, a basic walk through of import parts

## Application

### Application error handling

- when using `log.Fatal`, add a `return` after it, it won't get executed, but it makes the abort logic more obvious

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

## Library

### Library error handling

- DO NOT use `log.Fatal`, `panic`, always return error, if an error is added later, many application won't compile