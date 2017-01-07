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

func TestSet_Cardinality(t *testing.T) {
	assert := assert.New(t)
	s := NewSet("a", "b", "c")
	assert.Equal(3, s.Cardinality())
}

func TestSet_Size(t *testing.T) {
	assert := assert.New(t)
	s := NewSet("a", "b", "c")
	assert.Equal(s.Size(), s.Cardinality())
}

func TestSet_Equal(t *testing.T) {
	assert := assert.New(t)
	s := NewSet("a", "b", "c")
	s2 := NewSet("a")
	s3 := NewSet("a")
	assert.True(s.Equal(s))
	assert.False(s.Equal(s2))
	assert.True(s2.Equal(s3))
}
