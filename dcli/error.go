package dcli

import (
	"fmt"
	"reflect"
)

// error.go Defines special errors and handlers.

// ErrHelpOnlyCommand means the command does not contain execution logic and simply print usage info for sub commands.
type ErrHelpOnlyCommand struct {
	Command string
}

func NewErrHelpOnly(cmd string) *ErrHelpOnlyCommand {
	return &ErrHelpOnlyCommand{Command: cmd}
}

func (e *ErrHelpOnlyCommand) Error() string {
	return "command only prints help and has no execution logic: " + e.Command
}

func commandNotFound(root Command, args []string) {

}

// handles error from a specific command
// TODO: specify print output location for testing
func handleCommandError(cmd Command, err error) {
	switch x := err.(type) {
	case *ErrHelpOnlyCommand:
		fmt.Println("TODO: should print help for command " + cmd.GetName() + x.Command)
	// TODO: unwrap error? or by default simply print it ...
	default:
		fmt.Printf("TODO: unhandled error of type %s %s\n", reflect.TypeOf(err).String(), err)
	}
}
