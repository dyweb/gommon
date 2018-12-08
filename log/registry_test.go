package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_AddRegistry(t *testing.T) {
	rTop := Registry{}
	r1 := Registry{}
	rTop.AddRegistry(&r1)
	assert.Equal(t, 1, len(rTop.children))
	rTop.AddRegistry(&r1)
	assert.Equal(t, 1, len(rTop.children), "don't add registry if it already exists")
}

func TestRegistry_AddLogger(t *testing.T) {
	rTop := Registry{}
	l1 := NewTestLogger(InfoLevel)
	rTop.AddLogger(l1)
	assert.Equal(t, 1, len(rTop.loggers))
	rTop.AddLogger(l1)
	assert.Equal(t, 1, len(rTop.loggers), "don't add logger if it already exists")
}
