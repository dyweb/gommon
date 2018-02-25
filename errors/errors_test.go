package errors

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := asst.New(t)
	err := New("don't let me go")
	assert.NotNil(err)
	terr, ok := err.(TracedError)
	assert.True(ok)
	printFrames(terr.ErrorStack())
	// FIXME: this failed because frames can only be read once
	assert.Equal(3, framesLen(terr.ErrorStack()))
}
