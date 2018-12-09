# Changelog

## 2018-12-08

- remove logger relation ship from logger struct, use logger registry instead

## 2018-11-25

- reduce `Handler` interface method from 6 to 1, just a `HandleLog(level Level, now time.Time, msg string, source string, context Fields, fields Fields)`
- add `Print`, `Printf` for drop in replacement of standard library like non leveled logging library

The main takeaway for the benchmarks are 

- reduce number of alloc by allocate enough size in one call, it will reduce cpu time as well
- pass bytes slice all the way down the call stack to reduce alloc
- use interface methods cause parameters escape to heap, **even buffer pool can't help you with that** unless you eliminate interface use down the entire call stack
- you can still get good performance without using pool, logrus is using pool, but its performance is extremely poor