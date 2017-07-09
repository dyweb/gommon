# Config

Yaml configuration with template, inspired by [Ansible](http://docs.ansible.com/ansible/playbooks.html)

## YAML parsing problem

- multiple document support, the author said it is easy, but he never tried to solve them,
  - the [encoder & decoder PR](https://github.com/go-yaml/yaml/pull/163/) seems to support multiple documents but I don't
understand why it does.
  - Actually it would be interesting to look into how to unmarshal stream data since that is needed for Xephon-B,
though for most formats, go already have encoder and decoder
  - (Adopted) A easier way is to split the whole file by `---` before put it into go-yaml

## Specification

- only support single file, but can have multiple documents separated by `---`
  - pongo2 support include other templates from local file system, but we didn't make use of it
- variables
  - built in environment variable support, use `envs`
  - http://docs.ansible.com/ansible/playbooks_variables.html
- [ ] condition (not tested)
  - when
  - http://docs.ansible.com/ansible/playbooks_conditionals.html
- loop
  - http://docs.ansible.com/ansible/playbooks_loops.html
  - TODO: if using pongo2 syntax `for`, may have problem for order of rendering and parse yaml, unless support multiple
   document in one yaml file

Without using multiple document

- [ ] TODO: there is not function as `env()` shown in the example

````yaml
vars:
    influxdb_port: 8080
    databases:
        - influxdb
        - kairosdb
influxdb:
    port: "{{ influxdb_port }}"
kairosdb:
    port: "{{ env(\"KAIROSDB\") }}"
tests:
    - name: "test {{ item }}"
      cmd: "xephon-b load {{ item }} "
      with_items: "{{ databases }}"
    - name: "ping {{ item }}"
      cmd: "tsdb-proxy ping {{ item }} "
      with_items:
        - influxdb
        - kairosdb
````

Using multiple document

<!-- FIXED: it seems --- is treated as front matter http://assemble.io/docs/YAML-front-matter.html -->

````yaml

---
vars:
    influxdb_port: 8080
    databases:
        - influxdb
        - kairosdb
---
tests:
    (% for db in databases %)
    - name: "test {{ db }}"
      cmd: "xephon-b load {{ db }} "
   {% endfor %}
````

TODO: the example in runner actually does not require any template features

## Acknowledgement

- https://github.com/spf13/viper
- https://github.com/go-yaml/yaml
