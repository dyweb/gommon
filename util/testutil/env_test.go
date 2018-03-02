package testutil

import (
	"os"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestLoadDotEnv(t *testing.T) {
	assert := asst.New(t)
	LoadDotEnv(t)
	assert.Equal("BAR1", os.Getenv("FOO1"))
	assert.Equal("BAR2=BAR1", os.Getenv("FOO2"))
	assert.Equal("", os.Getenv("FOO3"))
}
