# Test utilities

TODO

- toc
- example in go doc, can there be example for test package?

## Condition

TODO: example of using test condition

## Env

TODO: load dot env

## Golden

TODO: use golden file

actually, just Read/WriteFixture would work ...

````go
// GenGolden check if env GOLDEN or GEN_GOLDEN is set, sometimes you need to generate test fixture in test
func GenGolden() Condition {
	return Or(EnvHas("GOLDEN"), EnvHas("GEN_GOLDEN"))
}
````