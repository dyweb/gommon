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
  
## Error in distributed tracing

- opencensus seems to be still working on it https://github.com/census-instrumentation/opencensus-proto/pull/145
  - I think I should ignore it for now ....

## Common self defined error types

- we can contains a helper func to convert error type to http status code
  - k8s is a example [reference 12-22](2018-12-22-error-types-reference.md)
- `AlreadyExists`
- `NotFound`
- `Unauthorized`
- `Forbidden`
- `NotImplemented`

## Standard library error

- `os.IsExist` `os.IsNotExist` `os.IsPermission`
- `os.PathError`
- `net.OpError`
- `net.Error`

## Reference

- [2018-12-14 Reference](2018-12-14-reference.md)