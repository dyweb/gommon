# 2018-09-21 Benchmark

There are three things, write the benchmark, run the benchmark, draw the graph. It's possible to collect profile as well,
i.e. use logger on a http server and load test the server, collect profile using pprof.

The main problem for putting benchmark in this repo is it will result in very large dependency, but if I don't put it
in this repo, then the develop would be very painful, need to copy unreleased file to the benchmark repo.

## Ref

Benchmarks

- zap https://github.com/uber-go/zap/tree/master/benchmarks
  - first clone and put the repo to $GOPATH/src/go.uber.org/zap  they are not using github repo as import path
- zerolog https://github.com/rs/logbench
- https://hackernoon.com/does-logging-cause-cpu-load-a-test-of-all-the-golang-logging-libraries-34052240f90d
  - [x] it measure system call using `sudo strace -c -t -p $(pid)`
- https://medium.com/justforfunc/analyzing-the-performance-of-go-functions-with-benchmarks-60b8162e61c6
- https://godoc.org/golang.org/x/perf

Tools

- https://godoc.org/golang.org/x/tools/benchmark/parse parse benchmark result