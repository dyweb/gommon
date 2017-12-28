# Seelog

https://github.com/cihub/seelog

- complexity: high

it has very detailed config (like log4j), also allow control logging behaviour of library.
besides level, it also support file and function pattern.

````xml
<seelog minlevel="info">
	<exceptions>
		<exception filepattern="test*" minlevel="error"/>
        <exception funcpattern="*test*" filepattern="tests.go" levels="off"/>
        <exception funcpattern="*perfCritical" minlevel="critical"/>
	</exceptions>
</seelog>
````

- patterns are matched using string comparison (not using regexp, custom recursive implementation)

````go
// logConfig stores logging configuration. Contains messages dispatcher, allowed log level rules
// (general constraints and exceptions)
type logConfig struct {
	Constraints    logLevelConstraints  // General log level rules (>min and <max, or set of allowed levels)
	Exceptions     []*LogLevelException // Exceptions to general rules for specific files or funcs
	RootDispatcher dispatcherInterface  // Root of output tree
}

// LogLevelException represents an exceptional case used when you need some specific files or funcs to
// override general constraints and to use their own.
type LogLevelException struct {
	funcPatternParts []string
	filePatternParts []string

	funcPattern string
	filePattern string

	constraints logLevelConstraints
}

// stringMatchesPattern check whether testString matches pattern with asterisks.
// Standard regexp functionality is not used here because of performance issues.
func stringMatchesPattern(patternparts []string, testString string) bool {
	// ...
}
````

````go
// innerLoggerInterface is an internal logging interface
type innerLoggerInterface interface {
	innerLog(level LogLevel, context LogContextInterface, message fmt.Stringer)
	Flush()
}

````

- log -> dispatcher -> write to output