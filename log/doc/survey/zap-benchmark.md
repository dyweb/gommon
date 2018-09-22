# Zap benchmark 

https://github.com/at15/zap/tree/master/benchmarks

- disabled level, no fields
  - zap has a check to avoid logging when it is enabled
- disabled level, fields attached to logger instance
- disabled level, fields attached when call
- log without field, pure message and use `f` for formatting
- log with fields attached to logger instance
- log with fields attached when call
  - [ ] TODO: what about combine fields, logger + call, it seems in zap, when you attach fields to logger, they are encoded right away

````text
⇒  go test -bench=.
goos: linux
goarch: amd64
pkg: go.uber.org/zap/benchmarks
BenchmarkDisabledWithoutFields/Zap-8            100000000               11.5 ns/op
BenchmarkDisabledWithoutFields/Zap.Check-8              100000000               11.4 ns/op
BenchmarkDisabledWithoutFields/Zap.Sugar-8              200000000                9.51 ns/op
BenchmarkDisabledWithoutFields/Zap.SugarFormatting-8            20000000                63.9 ns/op
BenchmarkDisabledWithoutFields/apex/log-8                       1000000000               2.59 ns/op
BenchmarkDisabledWithoutFields/sirupsen/logrus-8                200000000                7.62 ns/op
BenchmarkDisabledWithoutFields/rs/zerolog-8                     1000000000               2.40 ns/op
BenchmarkDisabledAccumulatedContext/Zap-8                       100000000               11.7 ns/op
BenchmarkDisabledAccumulatedContext/Zap.Check-8                 100000000               11.4 ns/op
BenchmarkDisabledAccumulatedContext/Zap.Sugar-8                 200000000                9.56 ns/op
BenchmarkDisabledAccumulatedContext/Zap.SugarFormatting-8       20000000                65.1 ns/op
BenchmarkDisabledAccumulatedContext/apex/log-8                  2000000000               1.06 ns/op
BenchmarkDisabledAccumulatedContext/sirupsen/logrus-8           200000000                7.24 ns/op
BenchmarkDisabledAccumulatedContext/rs/zerolog-8                1000000000               2.36 ns/op
BenchmarkDisabledAddingFields/Zap-8                             10000000               139 ns/op
BenchmarkDisabledAddingFields/Zap.Check-8                       100000000               11.6 ns/op
BenchmarkDisabledAddingFields/Zap.Sugar-8                       20000000                67.9 ns/op
BenchmarkDisabledAddingFields/apex/log-8                         5000000               242 ns/op
BenchmarkDisabledAddingFields/sirupsen/logrus-8                  3000000               421 ns/op
BenchmarkDisabledAddingFields/rs/zerolog-8                      30000000                47.9 ns/op
BenchmarkWithoutFields/Zap-8                                    10000000               147 ns/op
BenchmarkWithoutFields/Zap.Check-8                              10000000               146 ns/op
BenchmarkWithoutFields/Zap.CheckSampled-8                       30000000                44.0 ns/op
BenchmarkWithoutFields/Zap.Sugar-8                              10000000               222 ns/op
BenchmarkWithoutFields/Zap.SugarFormatting-8                      300000              3766 ns/op
BenchmarkWithoutFields/apex/log-8                                1000000              1992 ns/op
BenchmarkWithoutFields/go-kit/kit/log-8                          5000000               296 ns/op
BenchmarkWithoutFields/inconshreveable/log15-8                    500000              3279 ns/op
BenchmarkWithoutFields/sirupsen/logrus-8                         2000000               707 ns/op
BenchmarkWithoutFields/go.pedge.io/lion-8                        3000000               452 ns/op
BenchmarkWithoutFields/stdlib.Println-8                          3000000               427 ns/op
BenchmarkWithoutFields/stdlib.Printf-8                            500000              3162 ns/op
BenchmarkWithoutFields/rs/zerolog-8                             20000000               110 ns/op
BenchmarkWithoutFields/rs/zerolog.Formatting-8                    300000              4158 ns/op
BenchmarkWithoutFields/rs/zerolog.Check-8                       10000000               122 ns/op
BenchmarkAccumulatedContext/Zap-8                               10000000               158 ns/op
BenchmarkAccumulatedContext/Zap.Check-8                         10000000               157 ns/op
BenchmarkAccumulatedContext/Zap.CheckSampled-8                  30000000                45.7 ns/op
BenchmarkAccumulatedContext/Zap.Sugar-8                         10000000               233 ns/op
BenchmarkAccumulatedContext/Zap.SugarFormatting-8                 500000              3762 ns/op
BenchmarkAccumulatedContext/apex/log-8                             50000             26441 ns/op
BenchmarkAccumulatedContext/go-kit/kit/log-8                      200000              6196 ns/op
BenchmarkAccumulatedContext/inconshreveable/log15-8               100000             17316 ns/op
BenchmarkAccumulatedContext/sirupsen/logrus-8                     200000              7227 ns/op
BenchmarkAccumulatedContext/go.pedge.io/lion-8                    500000              2353 ns/op
BenchmarkAccumulatedContext/rs/zerolog-8                        20000000               113 ns/op
BenchmarkAccumulatedContext/rs/zerolog.Check-8                  20000000               114 ns/op
BenchmarkAccumulatedContext/rs/zerolog.Formatting-8               500000              3237 ns/op
BenchmarkAddingFields/Zap-8                                      1000000              1224 ns/op
BenchmarkAddingFields/Zap.Check-8                                1000000              1210 ns/op
BenchmarkAddingFields/Zap.CheckSampled-8                        10000000               166 ns/op
BenchmarkAddingFields/Zap.Sugar-8                                1000000              1448 ns/op
BenchmarkAddingFields/apex/log-8                                   50000             26937 ns/op
BenchmarkAddingFields/go-kit/kit/log-8                            200000              6039 ns/op
BenchmarkAddingFields/inconshreveable/log15-8                      50000             30591 ns/op
BenchmarkAddingFields/sirupsen/logrus-8                           200000              7216 ns/op
BenchmarkAddingFields/go.pedge.io/lion-8                          200000              5695 ns/op
BenchmarkAddingFields/rs/zerolog-8                                300000              5173 ns/op
BenchmarkAddingFields/rs/zerolog.Check-8                          300000              5158 ns/op
````