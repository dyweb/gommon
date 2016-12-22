package runner

import (
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

// commands that should be called using shell,
// because mostly they would be expecting shell expansion on parameters
// FIXME: use a map as set
var shellCommands = []string{
	"rm", "cp", "mv", "mkdir",
	"tar",
}

// NewCmd can properly split the executable with its arguments
// TODO: may need to add a context to handle things like using shell or not
func NewCmd(cmdStr string) (*exec.Cmd, error) {
	segments, err := shellquote.Split(cmdStr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse command")
	}
	return exec.Command(segments[0], segments[1:]...), nil
}

// NewCmdWithAutoShell automatically use `sh -c` syntax for a small list of executable
// because most people expect shell expansion i.e. wild chars when using them
// TODO: test
func NewCmdWithAutoShell(cmdStr string) (*exec.Cmd, error) {
	segments, err := shellquote.Split(cmdStr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse command")
	}
	name := segments[0]
	useShell := false
	for _, c := range shellCommands {
		if name == c {
			useShell = true
			break
		}
	}
	if useShell {
		// TODO: may use shellquote join?
		// NOTE: http://stackoverflow.com/questions/18946837/go-variadic-function-and-too-many-arguments
		// the `append` here is a must "sh", "-c", segments[1:]... won't work
		return exec.Command("sh", append([]string{"-c"}, segments[1:]...)...), nil
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
