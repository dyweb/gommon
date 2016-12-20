package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/dyweb/Ayi/issues/58
func TestNewCommand_ExtraQuote(t *testing.T) {
	assert := assert.New(t)
	cmd, _ := NewCmd("rm *.aux")
	assert.Equal(2, len(cmd.Args))
	// assert.Equal("/bin/rm", cmd.Path)
	assert.Equal("*.aux", cmd.Args[1])
	// TODO: it seems shell quote is right
}
