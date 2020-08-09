package dcli

import (
	"context"

	"github.com/dyweb/gommon/errors"
)

// command.go defines interface and the default implementation Cmd

type Runnable func(ctx context.Context) error

type Command interface {
	GetName() string
	GetRunnable() Runnable
	GetChildren() []Command
}

type MutableCommand interface {
	AddChildren(children ...Command) error
}

var _ Command = (*Cmd)(nil)

// Cmd is the default implementation of Command interface
type Cmd struct {
	Name     string
	Run      Runnable
	Children []Command
}

func (c *Cmd) GetName() string {
	return c.Name
}

func (c *Cmd) GetRunnable() Runnable {
	return c.Run
}

func (c *Cmd) GetChildren() []Command {
	return c.Children
}

func (c *Cmd) AddChildren(children ...Command) error {
	// TODO: check name conflict and import cycle
	c.Children = append(c.Children, children...)
	return nil
}

// Validate Start
const commandPrefixSep = ">"

// ValidateCommand checks if a command and its children have set name and runnable properly.
// It also checks if there is cycle TODO: the check is too strict ....
func ValidateCommand(c Command) error {
	m := make(map[Command]string)
	return validate(c, "", m)
}

func validate(c Command, prefix string, visited map[Command]string) error {
	merr := errors.NewMulti()
	name := "unknown"
	if c.GetName() == "" {
		merr.Append(errors.Errorf("command has no name, prefix: %s", prefix))
	} else {
		name = c.GetName()
	}
	if c.GetRunnable() == nil {
		merr.Append(errors.Errorf("command has no runnable, name: %s, prefix: %s", name, prefix))
	}
	prefix = prefix + commandPrefixSep + name
	// FIXME: this check is too strict, we only want cycle detection ... I think we allow DAG ...
	if p, ok := visited[c]; ok {
		merr.Append(errors.Errorf("duplicated command, previously used at %s used again at %s", p, prefix))
		return merr.ErrorOrNil()
	}
	visited[c] = prefix
	childNames := make(map[string]bool, len(c.GetChildren()))
	for _, child := range c.GetChildren() {
		if childNames[child.GetName()] {
			merr.Append(errors.Errorf("child defined twice, name: %s, parent: %s", child.GetName(), prefix))
		}
		merr.Append(validate(child, prefix, visited))
	}
	return merr.ErrorOrNil()
}

// Validate End

// Dispatch Start

func FindCommand(root Command, args []string) (Command, error) {
	if len(args) == 0 {
		return root, nil
	}
	// TODO: strip flag
	sub := args[0]
	for _, cur := range root.GetChildren() {
		if cur.GetName() != sub {
			continue
		}
		// Check if there can be more matches
		if len(cur.GetChildren()) == 0 || len(args) == 1 {
			return cur, nil
		}
		return FindCommand(cur, args[1:])
	}
	// TODO: typed error and suggestion using edit distance
	return nil, errors.New("command not found")
}

// Dispatch End

// Util Start

// hasChildCommand checks if command has child with given name. It does NOT check children recursively.
func hasChildCommand(c Command, name string) bool {
	for _, child := range c.GetChildren() {
		if child.GetName() == name {
			return true
		}
	}
	return false
}

// Util End
