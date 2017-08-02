package util

import (
	asst "github.com/stretchr/testify/assert"
	"testing"
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
