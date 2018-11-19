# 2018-11-18

Haven't worked on gommon for a while, thanks to the [go training](https://github.com/ardanlabs/gotraining), 
idea about go performance increased a bit. Especially what is allocation in go, previously I never really thought
about what is on heap and what is on stack

Continue on [2018-09-05](2018-09-05-clean-up.md) the basic steps are following

- finish benchmark regardless of results
  - try many log libraries in all the ways they allowed, structured, printf
    - zap
    - zerolog
    - glog?
    - stdlog
    - logrus
    - apex/log
  - if there are trivial optimization, optimize along the way (though the may be in vain once we decided to change the public interface)