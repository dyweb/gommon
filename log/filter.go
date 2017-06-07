package log

import (
	st "github.com/dyweb/gommon/structure"
)

// Filter determines if the entry should be logged
type Filter interface {
	Filter(entry *Entry) bool
	Name() string
}

// PkgFilter only allows entry without `pkg` field or `pkg` value in the allow set to pass
type PkgFilter struct {
	allow st.Set
}

// Filter checks if the pkg field is in the white list
func (filter *PkgFilter) Filter(entry *Entry) bool {
	pkg, ok := entry.Fields["pkg"]
	// entry without pkg is not filtered
	if !ok {
		return true
	}
	return filter.allow.Contains(pkg)
}

// Name implements Filter interface
func (filter *PkgFilter) Name() string {
	return "PkgFilter"
}

// NewPkgFilter returns a filter that allow log that contains `pkg` filed in the allow set
func NewPkgFilter(allow st.Set) *PkgFilter {
	return &PkgFilter{
		allow: allow,
	}
}
