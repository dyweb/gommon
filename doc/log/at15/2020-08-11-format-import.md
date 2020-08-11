# 2020-08-11 Format Import

## TODO

- [ ] format w/o considering comment
- [ ] format w/ comment

## Background

For [linter/import](../../../linter/doc/import.md). Most time is spent on reading goimport source.
The hard part is how to rearrange the AST and still generates a correct output w/ comment.
Go generates modified file by printing AST, and it also requires the modified `FileSet`.
It seems not all positions are adjusted properly, and they does not match the actual output in the end.
Removing blank line is easy (and used frequently). `fset.File(s.Pos()).MergeLine(fset.Position(p).Line)`
Adding blank line is harder, it is added after the code is already formatted using `addImportSpaces`

## Implementation

### Group

The grouping logic is same as `importGroup`

- define a set of rules, each rules gives a group number
- sort spec by group number like `byImportSpec`
- fix the comments and position like `sortSpecs`