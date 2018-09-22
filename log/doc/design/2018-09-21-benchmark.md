# 2018-09-21 Benchmark

There are three things, write the benchmark, run the benchmark, draw the graph. It's possible to collect profile as well,
i.e. use logger on a http server and load test the server, collect profile using pprof.

The main problem for putting benchmark in this repo is it will result in very large dependency, but if I don't put it
in this repo, then the develop would be very painful, need to copy unreleased file to the benchmark repo.