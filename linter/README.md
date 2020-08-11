# Gommon Linter

## Overview

Package `linter` is a code formatter and linter. It extends `goimports` with custom grouping rules.

## Install

- [ ] FIXME: rename the binary package to `gom` so we can download using go get

```bash
mkdir /tmp/gommon && cd /tmp/gommon && go get github.com/dyweb/gommon/cmd/gommon
```

## Usage

Format

- `gommon format` flags are compatible with `goimports` 

```bash
# print diff, list file and update file in place for go files under folder ./server ./client (recursively)
gommon format -d -l -w ./server ./client
```

## Internal

- [ ] TODO: talks about how the format works (maybe link to the blog post)
