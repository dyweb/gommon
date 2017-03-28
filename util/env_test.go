package util

import (
	"testing"
	"os"
	asst "github.com/stretchr/testify/assert"
)

func TestEnvAsMap(t *testing.T) {
	assert := asst.New(t)
	os.Setenv("gommondummy", "foo=bar")
	envMap := EnvAsMap()
	//t.Log(envMap)
	assert.Equal("foo=bar", envMap["gommondummy"])
}
