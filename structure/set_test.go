package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_Contains(t *testing.T) {
	assert := assert.New(t)
	s := NewSet("a", "b", "c")
	assert.True(s.Contains("a"))
	assert.False(s.Contains("d"))
}
