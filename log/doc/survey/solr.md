# Solr

The last straw for log v2 is Solr's admin UI, allow setting log level of different base

![solr-log-admin](solr-log-admin.png)

https://github.com/apache/lucene-solr

## Usage

- install [ant](http://ant.apache.org/bindownload.cgi)
- `mkdir -p ~/.ant/lib`
- `ant ivy-boostrap`
- `ant idea`

## Tree UI & Change level at runtime

- https://github.com/apache/lucene-solr/blob/master/solr/core/src/java/org/apache/solr/handler/admin/LoggingHandler.java

````java
if (set != null) {
    for (String pair : set) {
        String[] split = pair.split(":");
        String category = split[0];
        String level = split[1];
        watcher.setLogLevel(category, level);
    }
}
if (since != null) {
    // ... return recent log
} else {
    rsp.add("levels", watcher.getAllLevels());
    List<LoggerInfo> loggers = new ArrayList<>(watcher.getAllLoggers());
    Collections.sort(loggers);
    List<SimpleOrderedMap<?>> info = new ArrayList<>();
    for(LoggerInfo wrap:loggers) {
        info.add(wrap.getInfo());
    }
    rsp.add("loggers", info);
}
````

- https://github.com/apache/lucene-solr/tree/master/solr/core/src/java/org/apache/solr/logging
- https://github.com/apache/lucene-solr/blob/master/solr/core/src/java/org/apache/solr/logging/log4j/Log4jWatcher.java
  - it seems in this code, it split name by `.` to get the tree

````java
  @Override
  public Collection<LoggerInfo> getAllLoggers() {
    org.apache.log4j.Logger root = org.apache.log4j.LogManager.getRootLogger();
    Map<String,LoggerInfo> map = new HashMap<>();
    Enumeration<?> loggers = org.apache.log4j.LogManager.getCurrentLoggers();
    while (loggers.hasMoreElements()) {
      org.apache.log4j.Logger logger = (org.apache.log4j.Logger)loggers.nextElement();
      String name = logger.getName();
      if( logger == root) {
        continue;
      }
      map.put(name, new Log4jInfo(name, logger));

      while (true) {
        int dot = name.lastIndexOf(".");
        if (dot < 0)
          break;
        name = name.substring(0, dot);
        if(!map.containsKey(name)) {
          map.put(name, new Log4jInfo(name, null));
        }
      }
    }
    map.put(LoggerInfo.ROOT_NAME, new Log4jInfo(LoggerInfo.ROOT_NAME, root));
    return map.values();
  }
````

- https://github.com/apache/lucene-solr/blob/master/solr/webapp/web/js/angular/controllers/logging.js

````js
    var makeTree = function(loggers, packag) {
      var tree = [];
      for (var i=0; i<loggers.length; i++) {
        var logger = loggers[i];
        logger.packag = packageOf(logger);
        logger.short = shortNameOf(logger);
        if (logger.packag == packag) {
          logger.children = makeTree(loggers, logger.name);
          tree.push(logger);
        }
      }
      return tree;
    };
````

## Use Logger

- all the instances of a class is using same logger

````java
public class BlobHandler extends RequestHandlerBase implements PluginInfoInitialized , PermissionNameProvider {
    private static final Logger log = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());
}
````