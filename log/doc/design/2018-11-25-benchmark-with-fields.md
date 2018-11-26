# 2018-11-25 Benchmark with fields

Finally, gommon/log support fields in logger, which is often called context

## zap

Disabled

- disabled w/o fields
- disabled accumulated context (fields added to logger)
- disabled adding fields (adding fields when log)

Not disabled

- w/o fields
- accumulated context 
- adding fields

## zerolog

https://github.com/rs/logbench

It's more cleaner than zap's no copy and paste 

enable/disabled are top level

- no context
- with context
- different type of fields