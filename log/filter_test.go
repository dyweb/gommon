package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterInterface(t *testing.T) {
	t.Parallel()
	var _ Filter = (*PkgFilter)(nil)
}

func TestPkgFilter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	allow := make(map[string]bool)
	allow["ayi.app.git"] = true
	f := NewPkgFilter(allow)
	entryWithoutField := &Entry{}
	assert.True(f.Filter(entryWithoutField))
	field := make(map[string]string, 1)
	field["pkg"] = "ayi.app.git"
	entryWithAllowedPkg := &Entry{Fields: field}
	assert.True(f.Filter(entryWithAllowedPkg))
	field["pkg"] = "ayi.app.web"
	entryWithDisallowedPkg := &Entry{Fields: field}
	assert.False(f.Filter(entryWithDisallowedPkg))

}
