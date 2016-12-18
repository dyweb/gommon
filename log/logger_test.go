package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_AddFilter(t *testing.T) {
	assert := assert.New(t)
	logger := NewLogger()
	allow := make(map[string]bool)
	pkgFilter := NewPkgFilter(allow)
	logger.AddFilter(pkgFilter, DebugLevel)
	assert.Equal(1, len(logger.Filters[DebugLevel]))
}

func TestLogger_NewEntryWithPkg(t *testing.T) {
	assert := assert.New(t)
	logger := NewLogger()
	entry := logger.NewEntryWithPkg("x.dummy")
	assert.Equal(1, len(entry.Fields))
}