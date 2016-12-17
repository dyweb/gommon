package log

type Set map[string]bool

// Filter determines if the entry should be logged
type Filter interface {
	Filter(entry *Entry) bool
}

// PkgFilter only allows entry without `pkg` field or `pkg` value in the allow set to pass
type PkgFilter struct {
	allow Set
}

func (filter *PkgFilter) Filter(entry *Entry) bool {
	pkg, ok := entry.Fields["pkg"]
	// entry without pkg is not filtered
	if !ok {
		return true
	}
	_, ok = filter.allow[pkg]
	return ok
}

func NewPkgFilter(allow Set) *PkgFilter {
	return &PkgFilter{
		allow: allow,
	}
}
