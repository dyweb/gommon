package runner

import (
	"os/exec"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

// NewCmd can properly split the executable with its arguments
// TODO: may need to add a context to handle things like using shell or not
func NewCmd(cmdStr string) (*exec.Cmd, error) {
	segments, err := shellquote.Split(cmdStr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse command")
	}
	return exec.Command(segments[0], segments[1:]...), nil
}

func RunCommand(cmdStr string) error {
	cmd, err := NewCmd(cmdStr)
	if err != nil {
		return errors.Wrap(err, "can't create cmd from command string")
	}
	// TODO: dry run, maybe add a context
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failure when executing command")
	}
	return nil
}