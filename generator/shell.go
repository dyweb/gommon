package generator

import (
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
)

// https://github.com/dyweb/gommon/issues/53
type ShellConfig struct {
	Code  string `yaml:"code"`
	Shell bool   `yaml:"shell"`
	Cd    bool   `yaml:"cd"`
}

func (c *ShellConfig) Render(root string) error {
	log.Debugf("cmd %s shell %t cd %t", c.Code, c.Shell, c.Cd)
	var cmd *exec.Cmd
	if c.Shell {
		cmd = exec.Command("sh", "-c", c.Code)
	} else {
		if segments, err := shellquote.Split(c.Code); err != nil {
			return errors.Wrap(err, "can't split command into []string")
		} else {
			cmd = exec.Command(segments[0], segments[1:]...)
		}
	}
	if c.Cd {
		cmd.Dir = join(fsutil.Cwd(), root)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error executing command %s", c.Code)
	}
	return nil
}
