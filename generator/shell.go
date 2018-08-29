package generator

import (
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
)

// ShellConfig is shell command executed by gommon
// https://github.com/dyweb/gommon/issues/53
type ShellConfig struct {
	// Code is the command to be executed, Command will overwrite it if presented,
	// it is kept for backward compatibility
	Code    string `yaml:"code"`
	Command string `yaml:"command"`

	// Shell is true, use sh -c, otherwise use exec on the first segment after split
	Shell bool `yaml:"shell"`
	// Cd is true, switch to the folder of config file when executing command
	Cd bool `yaml:"cd"`
}

// Render execute the shell command, redirect stdin/out/err and block until it is finished
func (c *ShellConfig) Render(root string) error {
	command := c.Code
	if c.Command != "" {
		command = c.Command
	}
	log.Debugf("cmd %s shell %t cd %t", command, c.Shell, c.Cd)
	var cmd *exec.Cmd
	if c.Shell {
		cmd = exec.Command("sh", "-c", command)
	} else {
		if segments, err := shellquote.Split(command); err != nil {
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
