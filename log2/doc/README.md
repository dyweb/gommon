# gommon/log

## Survey

https://github.com/avelino/awesome-go#logging

- [std/log](std-log.md)
- [sirupsen/logrus](logrus.md)
- [uber-go/zap](zap.md)
- [apex/log](apex-log.md)
- [log15](log15.md)
- [gokit/log](gokit-log.md)
- [nanolog](nanolog.md)
- [glog](glog.md)
- [seelog](seelog.md)
- [log4j](log4j.md)
- [ ] TODO: might check open tracing as well, instrument like code should be put into other package

Logging library used by popular go projects

- k8s, [CockroachDB](https://github.com/cockroachdb/cockroach/tree/master/pkg/util/log) glog

## Goals

Major

- gives user fine grained control of logging in their application, including libraries that use gommon/log
  - like Java's, see Solr's admin page as example

![solr-log-admin](solr-log-admin.png)

- support level and context
- provide detail information when needed, i.e. file:line that can jump in Gogland

Minor

- util for filtering log data by package, file etc.
- benchmark
- test for race condition

Future

- sample, `LOG_IF`, `EVERY_N`, glog like flags etc.
- static code generation, i.e. expand `__FILE__`, `__LINE__`, `__FUNC__` instead of using `runtime`
- only log the delta like nanolog
- log in binary format

## Internals

Changes from v1

- `Entry` is now `Logger`, in v1, it contains both message (`Level`, `Time`, `Message`) and the real logging behaviour
  - it uses value receiver to make a copy when `log` because it updates `Level` etc.
  - first it calls the filters
  - then it calls `entry.Logger.Formatter.Format(&entry)` and pass itself
  - finally writes to output using `entry.Logger.Out.Write(serialized)`
- use `Handler`, but pass `Level`, `Time`, `Message`, `Fields` instead of a single `Entry` struct
  - handler takes care of both format and output
- `Fields` is now has two sorted slice, `[]string` for keys, `[]interface{}` for values
  - [ ] our use map? depends on how we want to use context, may not need to modify it
  - **two types of field**, one that don't change once set, one change in different log
    - i.e. protocol=http don't change, but code=500 may change to code=400