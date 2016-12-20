package runner

import (
	"os/exec"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

// NewCmd can properly split the executable with its arguments
func NewCmd(cmdStr string) (*exec.Cmd, error) {
	segments, err := shellquote.Split(cmdStr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse command")
	}
	return exec.Command(segments[0], segments[1:]...), nil
}
