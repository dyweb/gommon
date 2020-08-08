# Gommon Import Linter and Formatter

## TODO

- [ ] `LocalPrefix` kind of does the extra grouping
- [ ] Check forbidden import
- [ ] More groups for import

## Overview

`goimports` with custom grouping rules, package black list and alias naming check.

## Motivation

`goimports` improves `gofmt` by doing additional import grouping. However it is missing the following features:

- customize group rules, e.g. put all proto import at the bottom, split import from current project with external libs
- black list specific packages, sometimes IDE are too smart and introduced unknown import on the fly
- validation on import rename, e.g. rename lengthy proto package to `foopb`

## Design

`goimports` binary is a thin cli that calls `x/tools/imports` which calls `x/tools/internal/imports.Process`.
After merging, sorting and grouping import, it use `go/printer` to dump the ast and call `go/format.Source`.
(Essentially the code is pared twice, not sure why format.format is not exported).

It is possible to duplicate the format functionality in `x/tools/internal/imports`. However `goimports` can fix missing import,
and that logic is actually much larger than format and result in 17k lines of code for `imports` package.
So we take the easy way and run `goimport` before running `gommon/linter/import`.
One major drawback is the code is parsed and printed several times (in memory), so it should work for medium/small projects.

The overall flow is like following:

```
walkDir {
    src = readFile(p)
    b1 = imports.Process(p, src) // run goimports
    ast = parse(b1)
    err = lint.CheckImport(ast)
    b2 = lint.FormatImport(ast)
    diff(src, b2)
}
```

## Implementation

See [import.go](../import.go)

- `walkDir` and `diff` can copy from `goimports` binary until we have a better `fsutil.WalkWithIgnore` implementation
- `CheckImport` should use a default set of rules
  - so user can provide their own rules by writing their own binary
- `FormatImport` should use a default set of rules for grouping