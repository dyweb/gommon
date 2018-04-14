package testutil

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	asst "github.com/stretchr/testify/assert"
)

func TestCon_B(t *testing.T) {
	if IsTravis().B() {
		t.Log("you can use condition in if statement on travis")
	} else {
		t.Log("you can use condition in if statement")
	}
}

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

func TestAnd(t *testing.T) {
	assert := asst.New(t)

	res, msg, err := And(Bool("t", true), Bool("f", false)).Eval()
	assert.Nil(err)
	assert.False(res)
	assert.Equal("f is false", msg)
}

func TestOr(t *testing.T) {
	assert := asst.New(t)

	res, msg, err := Or(Bool("t", true), Bool("f", false)).Eval()
	assert.Nil(err)
	assert.True(res)
	assert.Equal("t is true", msg)
}

func TestNot(t *testing.T) {
	assert := asst.New(t)

	res, msg, err := Not(Bool("t", true)).Eval()
	assert.Nil(err)
	assert.False(res)
	assert.Equal("not t is true", msg)
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

func TestIsTravis(t *testing.T) {
	RunIf(t, IsTravis())
	t.Log("this is travis!")
}

func TestDump(t *testing.T) {
	if Dump().B() {
		spew.Dump(Dump())
	} else {
		t.Log("no dump")
	}
}
