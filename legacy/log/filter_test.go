// +build ignore

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
	// NOTE: we don't use package name with dot because
	// - when we get the package using reflection, we got /
	// - when access fields in config, we use dot notation, viper.Get("logging.ayi/app/git") is different with viper.Get("logging.ayi.app.gi")
	allow["ayi/app/git"] = true
	f := NewPkgFilter(allow)

	entryWithoutPkg := &Entry{}
	assert.False(f.Accept(entryWithoutPkg))
	entryWithAllowedPkg := &Entry{Pkg: "ayi/app/git"}
	assert.True(f.Accept(entryWithAllowedPkg))
	entryWithDisallowedPkg := &Entry{Pkg: "ayi/app/web"}
	assert.False(f.Accept(entryWithDisallowedPkg))

	// NOTE: we are using entry.Pkg instead of entry.Fields["pkg"]
	field := make(map[string]string, 1)
	field["pkg"] = "ayi/app/git"
	entryWithAllowedPkgInFields := &Entry{Fields: field}
	assert.False(f.Accept(entryWithAllowedPkgInFields))
}
