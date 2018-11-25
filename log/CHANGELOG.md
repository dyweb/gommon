# Changelog

## 2018-11-25

- reduce `Handler` interface method from 6 to 1, just a `HandleLog(level Level, now time.Time, msg string, source string, context Fields, fields Fields)`
- add `Print`, `Printf` for drop in replacement of standard library like non leveled logging library