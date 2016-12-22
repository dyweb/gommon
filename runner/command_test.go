package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/dyweb/Ayi/issues/58
// NOTE: it's not extra quote, it's os/exec can't expand * like shell does
// when run `rm *.aux` in shell, shell expands `*.aux` to real file names
func TestNewCommand_ExtraQuote(t *testing.T) {
	assert := assert.New(t)
	cmd, _ := NewCmd("rm *.aux")
	assert.Equal(2, len(cmd.Args))
	assert.Equal("*.aux", cmd.Args[1])
}

func TestNewCmdWithAutoShell(t *testing.T) {
	assert := assert.New(t)
	cmd, _ := NewCmdWithAutoShell("rm *.aux")
	assert.Equal(3, len(cmd.Args))
	// TODO: why this is not /bin/sh
	assert.Equal("sh", cmd.Args[0])
}