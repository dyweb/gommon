package log

import (
	st "github.com/dyweb/gommon/structure"
)

// Filter determines if the entry should be logged
type Filter interface {
	Accept(entry *Entry) bool
	FilterName() string
	FilterDescription() string
}

var _ Filter = (*PkgFilter)(nil)

// PkgFilter only allows entry without `pkg` field or `pkg` value in the allow set to pass
// TODO: we should support level
// TODO: a more efficient way might be trie tree and use `/` to divide package into segments instead of using character
type PkgFilter struct {
	allow st.Set
}

// Accept checks if the entry.Pkg (NOT entry.Fields["pkg"]) is in the white list
func (filter *PkgFilter) Accept(entry *Entry) bool {
	return filter.allow.Contains(entry.Pkg)
}

// FilterName implements Filter interface
func (filter *PkgFilter) FilterName() string {
	return "PkgFilter"
}

func (filter *PkgFilter) FilterDescription() string {
	return "Filter log based on their pkg tag value, it is logged if it does not have pkg field or in whitelist"
}

// NewPkgFilter returns a filter that allow log that contains `pkg` filed in the allow set
func NewPkgFilter(allow st.Set) *PkgFilter {
	return &PkgFilter{
		allow: allow,
	}
}
