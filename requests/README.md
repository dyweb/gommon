# Requests

A wrapper around net/http with less public global variables

Originally from [xephon-b](https://github.com/xephonhq/xephon-b)

## Usage

TODO: 

- it was mainly used for building config for http client to use socks5 proxy
- [ ] https://github.com/hashicorp/go-getter like features

## Changelog

- use same config as `DefaultTransport` in `net/http`
- add client struct and default client which wraps `http.Client`
- support socks5 proxy https://github.com/dyweb/gommon/issues/18
- rename `GetJSON` to get `GetJSONStringMap`, this should break old Xephon-B code, which is now in tsdb-proxy