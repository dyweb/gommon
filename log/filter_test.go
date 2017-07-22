package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestFilterInterface(t *testing.T) {
	t.Parallel()
	var _ Filter = (*PkgFilter)(nil)
}

func TestPkgFilter_Filter(t *testing.T) {
	t.Parallel()
	assert := asst.New(t)

	allow := make(map[string]bool)
	allow["ayi.app.git"] = true
	f := NewPkgFilter(allow)
	entryWithoutField := &Entry{}
	assert.True(f.Accept(entryWithoutField))
	field := make(map[string]string, 1)
	field["pkg"] = "ayi.app.git"
	entryWithAllowedPkg := &Entry{Fields: field}
	assert.True(f.Accept(entryWithAllowedPkg))
	field["pkg"] = "ayi.app.web"
	entryWithDisallowedPkg := &Entry{Fields: field}
	assert.False(f.Accept(entryWithDisallowedPkg))

}
