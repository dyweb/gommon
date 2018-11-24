# 2018-11-18

Haven't worked on gommon for a while, thanks to the [go training](https://github.com/ardanlabs/gotraining), 
idea about go performance increased a bit. Especially what is allocation in go, previously I never really thought
about what is on heap and what is on stack

Continue on [2018-09-05](2018-09-05-clean-up.md) the basic steps are following

- finish benchmark regardless of results
  - try many log libraries in all the ways they allowed, structured, printf, try their interface and see their output
    - zap
    - zerolog
    - glog?
    - stdlog
    - logrus
    - apex/log
  - if there are trivial optimization, optimize (or just keep track) along the way (though the may be in vain once we decided to change the public interface)
- list what are the use cases for gommon/log
- generate the modification plan
- execute plan
- another round of benchmark

## Using other logging libraries

- use fields to add context
- pass context down (i.e. store context in logger)
- source, file and line number

For adding context, there are two ways in existing library

- 'fast' libraries encode context to byte slice right away (result in they can't remove duplicated fields and sort fields added later)
- create an entry that holds both fields and reference to logger, attach all the logging methods `Debug` to the entry
  - logger's `Debug` is just `newEntry().Debug`

### Zap

Keep context (attach fields)

`logger.With(zap.Int("count", 1))`

- it will clone and return a new logger when adding new fields
  - also clone a core and encode fields right away into the core
  - core (similar to handler) **encode context data** (fields) to avoid encode it several times (exchange space for speed)
  - thus it will have duplicated fields if give same key with different data
  
### Zerolog

`logger.With().Str("k", "v").Logger()`

- when `Str` is called, it will encode field and append to `[]byte`

### Apex

`logger.WithField("f1", "v1").Info("have field")`

- same as logrus

### Logrus

`logger.WithField("f1", "v1").Info("have field")`

- create a new entry that hold pointer to logger, add fields to the new entry