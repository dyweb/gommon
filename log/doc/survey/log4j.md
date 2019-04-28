# Log4j

https://logging.apache.org/log4j/2.0/manual/api.html

Java 8 lambda support for lazy logging

````java
// pre-Java 8 style optimization: explicitly check the log level
// to make sure the expensiveOperation() method is only called if necessary
if (logger.isTraceEnabled()) {
    logger.trace("Some long-running operation returned {}", expensiveOperation());
}

// Java-8 lamdba
logger.trace("Some long-running operation returned {}", () -> expensiveOperation());
````

When looking at old [Solr survey](solr.md) and looking for its `getCurrentLoggers` method,
found it is [no longer the case in log4j 2](https://stackoverflow.com/a/18653927)

 ## V2 architecture

http://logging.apache.org/log4j/2.x/manual/architecture.html

Tl;DR configuration has the parent child relationship and updates on parent will be used in children

> In Log4j 1.x the Logger Hierarchy was maintained through a relationship between Loggers. In Log4j 2 this relationship no longer exists.
Instead, the hierarchy is maintained in the relationship between LoggerConfig objects

Logger names are case-sensitive and they follow the hierarchical naming rule

> A LoggerConfig is said to be an ancestor of another LoggerConfig if its name followed by a dot is a prefix of the descendant logger name. 
A LoggerConfig is said to be a parent of a child LoggerConfig if there are no ancestors between itself and the descendant LoggerConfig

- logger has name, and the best known strategy so far is using fully qualified class names (fqcn ...)
- logger config can reference its parent

Filter

- called before appender
- three results
  - accept
  - deny
  - neutral, pass to other filter

Appender

- allow log to multiple location, console, file, remote socket ....

> An Appender can be added to a Logger by calling the addLoggerAppender method of the current Configuration. 
If a LoggerConfig matching the name of the Logger does not exist, one will be created, 
the Appender will be attached to it and then all Loggers will be notified to update their LoggerConfig references

- **all Loggers will be notified to update their LoggerConfig references**

> Each enabled logging request for a given logger will be forwarded to all the appenders in that Logger's LoggerConfig as well as the Appenders of the LoggerConfig's parents. 
In other words, Appenders are inherited additively from the LoggerConfig hierarchy. 