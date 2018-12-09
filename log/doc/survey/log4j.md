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