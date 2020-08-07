# Go Log Library Survey

https://github.com/avelino/awesome-go#logging

High performance

- [uber-go/zap](zap.md) structured logging with high performance
- [rs/zerolog](zerolog.md) high performance, only focus on json logging
- [nanolog](nanolog.md) use binary format, only store the changed part (like prepare statement in database)

Simple

- [std/log](std-log.md) standard library log package
- [glog](glog.md) leveled only, but can sample based on number of hits on a certain file:line
- [gokit/log](gokit-log.md) extreme simple interface

Structured

- [sirupsen/logrus](logrus.md) structured logging, poor performance
- [apex/log](apex-log.md) use handler instead of formatter + writer
- [log15](log15.md) lazy evaluation

Java(ish)

- [solr](solr.md) the last straw that drives us to log v2, gives you [a tree graph to control log level of ALL the packages](solr-log-admin.png), including dependencies
- [seelog](seelog.md) javaish, fine grained control log filtering (by func, file etc.) at log site
- [log4j](log4j.md) java logger
- [ ] TODO: might check open tracing as well, instrument like code should be put into other package

Logging library used by popular go projects

- k8s, [CockroachDB](https://github.com/cockroachdb/cockroach/tree/master/pkg/util/log) glog

Rotate

- [ ] https://github.com/lestrrat-go/file-rotatelogs