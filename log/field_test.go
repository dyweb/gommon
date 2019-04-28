package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFields(t *testing.T) {
	fields := Fields{
		Str("k1", "v1"),
		Int("k2", 2),
	}
	copied := copyFields(fields)
	assert.Equal(t, 2, cap(copied), "capacity is same as length")
	fields[0].Key = "k1modified"
	copied[0].Str = "v1modified"

	assert.Equal(t, "k1modified", fields[0].Key)
	assert.Equal(t, "v1", fields[0].Str)
	assert.Equal(t, "k1", copied[0].Key)
	assert.Equal(t, "v1modified", copied[0].Str)
}
