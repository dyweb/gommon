package dcli_test

import (
	"testing"

	"github.com/dyweb/gommon/dcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCommand(t *testing.T) {
	// TODO: flag, sub command, auto complete suggestion
	r := &dcli.Cmd{
		Name: "bh",
		Children: []dcli.Command{
			&dcli.Cmd{
				Name: "user",
				Children: []dcli.Command{
					&dcli.Cmd{
						Name: "register",
					},
				},
			},
		},
	}
	c, err := dcli.FindCommand(r, []string{"user", "register"})
	require.Nil(t, err)
	assert.Equal(t, "register", c.GetName())
}
