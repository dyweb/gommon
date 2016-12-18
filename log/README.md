# Log that support filter by field

This a simple re-implementation of [logrus](https://github.com/sirupsen/logrus)

## Added

- `Filter` and `PkgFilter`
- `TraceLevel`

## Fixed

- Remove duplicate code in `Logger` and `Entry` by only allow using `Entry` to log
- Use elasped time in `TextFormatter`, see [issue](https://github.com/sirupsen/logrus/issues/457)
- `FatalLevel` should be more severe than `PanicLevel`

## Removed 

- lock on logger when call log for `Entry`
- Support for blocking Hook
- Trim `\n` when using `*f`

## TODO

- [ ] read filters from command line or config files
- [ ] `WithFields`
- [ ] pool for `Entry` and `bytes.Writer`
- [ ] JSON Formatter
- [ ] Multiple output
- [ ] async Hook
- [ ] Batch write and flush like [zap](https://github.com/uber-go/zap)
- [ ] Shutdown handler

## DONE

- thread safe by not using pointer receiver https://github.com/at15/go-learning/issues/3
- log to stdout, implemented `TextFormatter`
- leveled logging, add `Trace` level
- support field
- filter by field