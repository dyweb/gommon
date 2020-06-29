# Attribution & Comparison

Gommon is inspired by many awesome libraries. However, we chose to reinvent the wheel for most functionalities.
Doing so allow us to introduce break changes frequently ...

## errors

- [pkg/errors](https://github.com/pkg/errors) it cannot introduce breaking change, but `WithMessage` and `WithStack` is annoying
  - see [#54](https://github.com/dyweb/gommon/issues/54) and [errors/doc](../errors/doc) about other error packages
  - https://github.com/pkg/errors/pull/122 implemented checking existing stack before attach new one
- [uber-go/multierr#21]( https://github.com/uber-go/multierr/issues/21) for return bool after append
- [hashicorp/go-multierror](https://github.com/hashicorp/go-multierror) for `ErrorOrNil`

## log

- [sirupsen/logrus](https://github.com/sirupsen/logrus) for structured logging 
  - log v1 is entirely modeled after logrus, entry contains log information with methods like `Info`, `Infof`
- [apex/log](https://github.com/apex/log) for log handlers
  - log v2's handler is inspired by apex/log, but we didn't use entry struct and chose to pass multiple parameters to explicitly state what a handler should handle.
  This can avoid the problem of adding extra context and ignored silently by existing handler implementations.
- [uber-go/zap](https://github.com/uber-go/zap) for serialize log fields without using `fmt.Sprintf` and use `strconv` directly
  - we didn't go that extreme as Zap or ZeroLog for zero allocation, performance is currently not measured

## noodle

- [GeertJohan/go.rice](https://github.com/GeertJohan/go.rice)
  - we implemented `.gitignore` like [feature](https://github.com/at15/go.rice/issues/1) but the upstream didn't respond for the [feature request #83](https://github.com/GeertJohan/go.rice/issues/83)
  - we put data into generated code file while go.rice and append zip to existing go binary

## config

- [spf13/cast](https://github.com/spf13/cast) for cast, it is used by Viper
- [spf13/viper](https://github.com/spf13/viper/) for config
  - looking up config via string key makes type system useless, so we always marshal entire config file to a single struct
    - it also makes refactor easier
  - we only use YAML, might add RCL

## requests

- [python requests](http://docs.python-requests.org/en/master/) for requests
- [hashicorp/go-cleanhttp](https://github.com/hashicorp/go-cleanhttp) for using non default http transport and client

## generator

- [benbjohnson/tmpl](https://github.com/benbjohnson/tmpl) for go template based generator
  - first saw it in [influxdata/influxdb](https://github.com/influxdata/influxdb/blob/master/tsdb/engine/tsm1/encoding.gen.go.tmpl)
  - we put template data in `gommon.yml`, so we don't need to pass data as json via cli.
  - Using YAML instead of flags based on [docker-compose](https://github.com/docker/compose)

