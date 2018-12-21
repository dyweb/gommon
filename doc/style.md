# Coding Style

General Go coding style, this should be followed by people developing gommon and [other Go projects
under dyweb](https://github.com/dyweb?utf8=%E2%9C%93&q=&type=&language=go)

## Package

Based on [Go Blog: Package names][Go: Package names]

- use small packages instead of single flat package
  - a large package can have many import and is easier to cause import cycle if it is a library
  - it's easy to merge smaller packages into large one but hard in the other direction (TODO: refer hashicorp blog)
- package name should be same as folder name
  - avoid `go-mylib`, `go.mylib` as end of import path for `mylib` (we know it's go code ...)
- use `alllowercase` no `camelCase` or `snake_case`
  - `camelCase` cause trouble for some filesystems, i.e. windows users (even worse when they mount it into a linux vm)
  - `snake_case` normally have special meaning in go, i.e. `foo_test.go` (it's test file), `fs_windows.go` (only build for windows platform)
  - this `hardtoreadname` also distinguish it from other identifiers like `varName` `funcName`
- rename your package if all the usage have to rename it
  - when user rename the import and copy paste the code to another file, editors can't fill in the right import
  - people can have many different renames even inside a single project `httpUtil`, `httpUtility`, 
  it's better just give people a name that don't conflict and can use it everywhere
- use `util.go` instead of `util` package
  - a lot of times, those `util` are only used once in current package
  - you will have a hard time rename all the `util` packages
  
## File

- put definitions on the top
  - constant, package level vars (sentiel errors)
  - exported interface
  - exported struct
- group implementations together
  - if a interface have multiple implementations, put them close to each other, but **do NOT intersect the methods from
  different structs**, this only makes copy and paste easier
- keep file small
  - if a file have many struct and many methods, you should group them by functionality and split into smaller files
  - if a struct has too many/lengthy methods in one file, you should consider decouple the struct

## Struct

- use struct literal initialization
  - you type less and rename variable is easier when there is no editor's help
- do NOT use embedding for struct that will be serialized (in JSON etc.)
- prefer function over struct methods
  - TODO: explain this ...

## Func

- when pass/return a pointer, use `&` at call site, NOT initialization
  - pointer means ownership
  - pointer does not means no copy, sometimes it causes allocation on heap which can be avoided
  - [ ] TODO: refer to Bill's blog and course
   
## Test

- use helper function that accepts `testing.T` and don't return error
  - your code should handle error gracefully, your test should stop when there is error/unexpected result
  - this may cause error source line harder to find since go's test fail function does allow pass skip
  - [ ] TODO: might open a issue in go's repo, or there are already issues like that
- use [subtest][Go: Subtest and Sub-benchmarks]
  - this allows setup and teardown using vanilla go code without any framework
- use `TestMain` for package level setup and tear down
- [ ] TODO: add godlen file, add refer to hashicorp test video

## Ref

[Go: Package names]: https://blog.golang.org/package-names
[Go: Subtest and Sub-benchmarks]: https://blog.golang.org/subtests