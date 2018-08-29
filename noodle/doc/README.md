# Noodle

## TODO

- [x] generate go file based on assets (like most bind-data ish library does)
- [ ] figure out how to create and append zipfile, without relying on external zip binary

## Implementation

Current implementation generate a go file with a huge `[]byte` and metadata as go struct, all the files are compressed
into single zip file, folders are only kept as go struct.

## Ref

- https://github.com/GeertJohan/go.rice 1.6k star, used to use it for Ayi
- https://github.com/shurcooL/vfsgen 236 star listed most packages in alternative section
  - https://github.com/shurcooL/vfsgen#alternatives
- https://github.com/rakyll/statik 1,4k star
- https://www.ueber.net/who/mjl/blog/p/assets-in-go-binaries/ describes how to append zip
- https://github.com/cookieo9/resources-go pretty like go.rice
- [golang-nuts/including html contents of a file at compile time](https://groups.google.com/forum/#!topic/golang-nuts/9QUjmDED96E)