# Benchmark

## Log libraries

- [ ] TODO: k8s fork for glog https://github.com/kubernetes/klog , 
they are also considering parent children logger https://github.com/kubernetes/klog/issues/22

## Ref

- zap https://github.com/uber-go/zap/tree/master/benchmarks
  - first clone and put the repo to $GOPATH/src/go.uber.org/zap  they are not using github repo as import path
- zerolog https://github.com/rs/logbench
- https://hackernoon.com/does-logging-cause-cpu-load-a-test-of-all-the-golang-logging-libraries-34052240f90d
  - it starts a http server to log and use a client to hit the server
  - [x] it measure system call using `sudo strace -c -t -p $(pid)` and see context switches
- https://medium.com/justforfunc/analyzing-the-performance-of-go-functions-with-benchmarks-60b8162e61c6