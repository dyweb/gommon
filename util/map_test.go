package util

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestMergeStringMap(t *testing.T) {
	assert := asst.New(t)
	m := make(map[string]interface{})
	m2 := make(map[string]interface{})
	m["a"] = 123
	m2["a"] = 124
	m2["b"] = 125
	MergeStringMap(m, m2)
	assert.Equal(124, m["a"])
	assert.Equal(2, len(m))
}

func TestMapKeys(t *testing.T) {
	assert := asst.New(t)
	// nil
	assert.Empty(MapKeys(nil))
	assert.Empty(MapSortedKeys(nil, true))

	// empty map
	m := make(map[string]interface{})
	assert.Empty(MapKeys(m))
	assert.Empty(MapSortedKeys(m, true))

	m["a"] = 123
	m["b"] = 124
	assert.Equal([]string{"a", "b"}, MapSortedKeys(m, true))
	assert.Equal([]string{"b", "a"}, MapSortedKeys(m, false))
}
