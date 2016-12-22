package log

import (
	st "github.com/dyweb/Ayi/common/structure"
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

func (filter *PkgFilter) Filter(entry *Entry) bool {
	pkg, ok := entry.Fields["pkg"]
	// entry without pkg is not filtered
	if !ok {
		return true
	}
	return filter.allow.Contains(pkg)
}

func (filter *PkgFilter) Name() string {
	return "PkgFilter"
}

func NewPkgFilter(allow st.Set) *PkgFilter {
	return &PkgFilter{
		allow: allow,
	}
}
