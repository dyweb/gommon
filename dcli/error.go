package dcli

import (
	"fmt"
	"reflect"
)

// error.go Defines special errors and handlers.

// ErrHelpOnlyCommand means the command does not contain execution logic and simply print usage info for sub commands.
// TODO: we should consider add a factory method
type ErrHelpOnlyCommand struct {
	Command string
}

func (e *ErrHelpOnlyCommand) Error() string {
	return "command only prints help and has no execution logic: " + e.Command
}

func commandNotFound(root Command, args []string) {

}

// handles error from a specific command
func handleCommandError(cmd Command, err error) {
	switch x := err.(type) {
	case *ErrHelpOnlyCommand:
		fmt.Println("TODO: should print help for command " + cmd.GetName() + x.Command)
	default:
		fmt.Printf("TODO: unhandled error of type %s", reflect.TypeOf(err).String())
	}
}
