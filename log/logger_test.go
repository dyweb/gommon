package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFilter(t *testing.T) {
	assert := assert.New(t)
	logger := NewLogger()
	allow := make(map[string]bool)
	pkgFilter := NewPkgFilter(allow)
	logger.AddFilter(pkgFilter, DebugLevel)
	assert.Equal(1, len(logger.Filters[DebugLevel]))
}
