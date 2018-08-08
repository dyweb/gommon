# Config

Yaml configuration with template, inspired by [Ansible](http://docs.ansible.com/ansible/playbooks.html) and [Viper](https://github.com/spf13/viper)

## Usage

- [ ] TODO

## Specification

- using [golang's template](https://golang.org/pkg/text/template/) syntax, not ansible, not pongo2
- only support single file, but can have multiple documents separated by `---`
- environment variables
  - use `env` function, i.e. `home: {{ env "HOME"}}`
- variables
  - special top level key `vars` is used via `var` function
  - `vars` in different documents are merged, while for other top level keys, their value is replaced by latter documents
  - `vars` in current document can be used in current document, we actually render the template of each document twice, so in vars
  section you can use previous `vars` and all the syntax supported by golang's template
- [ ] condition (not tested, but golang template should support it)
- loop
  - use the `range` syntax

Using variables

- [ ] TODO: will `xkb --target={{ $name }} --port={{ $db.port }}` have undesired blank and/or new line?

````yaml
vars:
    influxdb_port: 8080
    databases:
        - influxdb
        - kairosdb
influxdb:
    port: {{ var "influxdb_port" }}
kairosdb:
    port: {{ env "KAIROSDB_PORT" }}
tests:
{{ range $name := var "databases" -}}
{{ $db := var $name -}}
    - xkb --target={{ $name }} --port={{ $db.port }}
{{ end }
base: {{ env "HOME" }}
````

Using multiple document

<!-- FIXED: it seems --- is treated as front matter http://assemble.io/docs/YAML-front-matter.html -->

- [ ] TODO: example
- [ ] TODO: the example in runner actually does not require any template features

## YAML parsing problem

- multiple document support, the author said it is easy, but he never tried to solve them,
  - the [encoder & decoder PR](https://github.com/go-yaml/yaml/pull/163/) seems to support multiple documents but I don't
understand why it does.
  - Actually it would be interesting to look into how to unmarshal stream data since that is needed for Xephon-B,
though for most formats, go already have encoder and decoder
  - (Adopted) A easier way is to split the whole file by `---` before put it into go-yaml
  
## Acknowledgement

- https://github.com/spf13/viper
- https://github.com/go-yaml/yaml
