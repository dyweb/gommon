package errors

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestStack_Frames(t *testing.T) {
	assert := asst.New(t)
	s := callers()
	// when Frames() is not called, it is empty
	assert.Equal(0, len(s.frames))
	PrintFrames(s.Frames())
	assert.Equal(s.depth, len(s.Frames()))
}
