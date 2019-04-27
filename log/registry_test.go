package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_AddLogger(t *testing.T) {
	rTop := Registry{}
	l1 := NewTestLogger(InfoLevel)
	rTop.addLogger(l1)
	assert.Equal(t, 1, len(rTop.loggers))
	rTop.addLogger(l1)
	assert.Equal(t, 1, len(rTop.loggers), "don't add logger if it already exists")
}
