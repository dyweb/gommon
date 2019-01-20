# 2019-01-02 Error types

This doc describes the goals for errors package in v0.0.10, need to continue the original goal for v0.0.9.

- `errortypes` package with `IsXXX` wrapper on built in errors and interface for common application errors
  - package level exported reflect.Value for caching
  - common application error types, see [reference 12-14](2018-12-14-reference.md) and
[reference 12-22](2018-12-22-error-types-reference.md)
  - [ ] might need to consider validation errors, could put it in another package and define interface in error types
- (optional) error collecting interface
  - could be a simple in memory service and have `errors/repoter.Save(err)` and later send to remote
  - a lot of (distributed) tracing systems are taking error collection into account as well
  
## Reference

- [2018-12-14 Reference](2018-12-14-reference.md)