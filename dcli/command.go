package dcli

import "context"

type Command interface {
	Name() string
	Run(ctx context.Context) error
}

type Cmd struct {
}
