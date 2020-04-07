package dcli

import (
	"context"

	"github.com/dyweb/gommon/errors"
)

type Runner func(ctx context.Context) error

type Command interface {
	GetName() string
	GetRun() Runner
	GetChildren() []Command
}

// validate checks if all the commands have set name and runnable properly.
func validate(c Command) error {
	merr := errors.NewMultiErr()
	if c.GetName() == "" {
		merr.Append(errors.New("command has no name"))
	}
	if c.GetRun() == nil {
		merr.Append(errors.Errorf("command %s has no runner", c.GetName()))
	}
	for _, child := range c.GetChildren() {
		merr.Append(validate(child))
	}
	return merr.ErrorOrNil()
}

// Cmd is the default implementation of Command interface
type Cmd struct {
	Name     string
	Run      Runner
	Children []Command
}

func (c *Cmd) GetName() string {
	return c.Name
}

func (c *Cmd) GetRun() Runner {
	return c.Run
}

func (c *Cmd) GetChildren() []Command {
	return c.Children
}
