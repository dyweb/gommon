# Requests

A pythonic HTTP library for Gopher with socks5 proxy support

Originally from [xephon-b](https://github.com/xephonhq/xephon-b)

## Usage

## Changelog

- add client struct and default client which wraps `http.Client`
- support socks5 proxy https://github.com/dyweb/gommon/issues/18
- rename `GetJSON` to get `GetJSONStringMap`, this should break old Xephon-B code, which is now in tsdb-proxy