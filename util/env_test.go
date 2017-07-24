package util

import (
	"testing"
	"os"
	asst "github.com/stretchr/testify/assert"
)

func TestLoadDotEnv(t *testing.T) {
	assert := asst.New(t)
	LoadDotEnv(t)
	assert.Equal("BAR1", os.Getenv("FOO1"))
	assert.Equal("BAR2=BAR1", os.Getenv("FOO2"))
}

func TestEnvAsMap(t *testing.T) {
	assert := asst.New(t)
	os.Setenv("gommondummy", "foo=bar")
	envMap := EnvAsMap()
	//t.Log(envMap)
	assert.Equal("foo=bar", envMap["gommondummy"])
}
