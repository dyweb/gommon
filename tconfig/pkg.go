// Package tconfig is a traceable config package. It allows you to keep track of the source and mutation
// of single or a set of config values from different sources, e.g. environment variable, config file, cli flags.
package tconfig

type Var interface {
	Eval() interface{}
}

type BoolVar interface {
	EvalBool() bool
}

type IntVar interface {
	EvalInt() int
}

type StringVar interface {
	EvalString() string
}
