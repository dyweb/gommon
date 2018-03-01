package testutil

import (
	"testing"
	"os"

	asst "github.com/stretchr/testify/assert"
)

func TestEnvHas(t *testing.T) {
	assert := asst.New(t)

	res, msg, err := EnvHas("USER").Eval()
	assert.True(res)
	assert.Equal("env USER is defined", msg)
	assert.Nil(err)

	res, msg, err = EnvHas("HOUXIYA_MINI_YOUKOU").Eval()
	assert.False(res)
	assert.Nil(err)
}

func TestEnvTrue(t *testing.T) {
	assert := asst.New(t)

	err := os.Setenv("THIS_MUST_BE_TRUE", "1")
	assert.Nil(err)
	res, msg, err := EnvTrue("THIS_MUST_BE_TRUE").Eval()
	assert.True(res)
	assert.Equal("env THIS_MUST_BE_TRUE=true", msg)
	assert.Nil(err)
}

func TestRunIf(t *testing.T) {
	t.Run("travis", func(t *testing.T) {
		RunIf(t, EnvTrue("TRAVIS"))
		t.Log("I should be seen on travis only!")
	})
}

func TestSkipIf(t *testing.T) {
	t.Run("travis", func(t *testing.T) {
		SkipIf(t, EnvTrue("TRAVIS"))
		t.Log("I should NOT be seen on travis")
	})
}
