# 2020-08-06 Pending Issues

## Background

I decided to follow the project management in [BenchHub](https://github.com/benchhub/benchhub).
i.e. moving project planning from github issue to markdown file in milestones and components.
There are many pending issues across a long time span and I need to summarize them before drafting new plan.

The [issues](https://github.com/dyweb/gommon/issues?page=1&q=is%3Aissue+is%3Aopen) can be divided into the following categories:

- new package
- new features on existing packages
- bug

Some issues are so old the original packages are removed (config, requests etc.).

## Issues

### New package

Large

- [dcli](https://github.com/dyweb/gommon/issues/117) a cli builder to replace cobra
- [tail](https://github.com/dyweb/gommon/issues/95) like `tail -f` and might even parse log and send metrics to tsdb like [mtail](https://github.com/google/mtail)
- [testx](https://github.com/dyweb/gommon/issues/101) test and benchmark result dashboard

Medium

- [goimport that checks and replaces specific import](https://github.com/dyweb/gommon/issues/118)

Small

- [human](https://github.com/dyweb/gommon/issues/10)
- [retry](https://github.com/dyweb/gommon/issues/126)
- [mathutil](https://github.com/dyweb/gommon/issues/123) `mathuitl.MaxInt(64)`
- [netutil](https://github.com/dyweb/gommon/issues/122) port wait for it

### New feature

- error
  - [a more complex error interface](https://github.com/dyweb/gommon/issues/76)
  - [human readable suggestion for possible solutions](https://github.com/dyweb/gommon/issues/73)
  - [fmt.Formatter](https://github.com/dyweb/gommon/issues/62) I think the new go error package has abandoned this
- log
  - [rename log to dlog](https://github.com/dyweb/gommon/issues/120)
  - [parse generated log](https://github.com/dyweb/gommon/issues/89)
  - [generate file and line number to avoid calling runtime](https://github.com/dyweb/gommon/issues/43)
  - [http API to control log level at runtime](https://github.com/dyweb/gommon/issues/23)
  - [support grep log and web UI](https://github.com/dyweb/gommon/issues/9)
- noodle
  - [set modification time for generated file](https://github.com/dyweb/gommon/issues/128)
  - [interface around http.FileSystem](https://github.com/dyweb/gommon/issues/84)
- generator
  - [deepcopy](https://github.com/dyweb/gommon/issues/102)
- testutil
  - [Only run test in IDE](https://github.com/dyweb/gommon/issues/91)
- requests (the package is in legacy already and replaced by httpclient)
  - [oauth2 client with access token](https://github.com/dyweb/gommon/issues/70)
- config
  - [validate config struct using tag](https://github.com/dyweb/gommon/issues/19)

### Bug

- [go vet error in example](https://github.com/dyweb/gommon/issues/107)

## Priority

It's impossible to fix all the gommon issues at once, and there are several active projects using gommon (benchhub, gce4-go, pm).
Some features are nice to have e.g. adjust log level using http API while some features are essential e.g. dcli.

- dcli
  - used by all projects that requires a cli, benchhub, gce4-go, pm, ayi
- log
  - used by all projects, the API is hard to use, and it's hard to read the log
- error
  - used by all projects, but most time I am just wrapping w/o analysing i.e. no unwrap
- generator
  - used by projects that uses protobuf, benchhub
- test
  - helps develop all the go packages
  - used by benchhub for gobench and gotest framework
- util
  - mathutil, stringutil, maputil etc.
- noodle
  - I rarely use it because I read from local fs, it will be useful for projects that distribute binary w/ UI