package util

import (
	asst "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	assert := asst.New(t)
	LoadDotEnv(t)
	assert.Equal("BAR1", os.Getenv("FOO1"))
	assert.Equal("BAR2=BAR1", os.Getenv("FOO2"))
	assert.Equal("", os.Getenv("FOO3"))
}

func TestEnvAsMap(t *testing.T) {
	assert := asst.New(t)
	os.Setenv("gommondummy", "foo=bar")
	envMap := EnvAsMap()
	//t.Log(envMap)
	assert.Equal("foo=bar", envMap["gommondummy"])
}
