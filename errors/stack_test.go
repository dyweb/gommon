package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_Frames(t *testing.T) {
	s := callers()
	// when Frames() is not called, it is empty
	assert.Equal(t, 0, len(s.frames))
	PrintFrames(s.Frames())
	assert.Equal(t, s.depth, len(s.Frames()))
}
