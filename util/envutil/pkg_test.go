package envutil

import (
	"os"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

// TODO: test LoadDotEnv

func TestEnvAsMap(t *testing.T) {
	assert := asst.New(t)
	os.Setenv("gommondummy", "foo=bar")
	envMap := EnvMap()
	//t.Log(envMap)
	assert.Equal("foo=bar", envMap["gommondummy"])
}
