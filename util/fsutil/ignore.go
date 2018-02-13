package fsutil

import (
	"path/filepath"
)

type Ignores []IgnorePattern

func (is *Ignores) ShouldIgnore(path string) bool {
	for _, p := range *is {
		if p.ShouldIgnore(path) {
			return true
		}
	}
	return false
}

type IgnorePattern interface {
	ShouldIgnore(path string) bool
}

type ExactPattern string

func (p ExactPattern) ShouldIgnore(path string) bool {
	return string(p) == path
}

// Deprecated it is not implemented yet
// NOTE: only * and ? is supported
// * matches any non empty sequence of non-separator character
// ? matches one non-separator character
// we are NOT expecting path to have * and ?
type WildcardPattern string

func (p WildcardPattern) ShouldIgnore(path string) bool {
	// the pattern would always be no greater than path due to our limited features
	if len(p) > len(path) {
		return false
	}
	s := -1 // * location
	i := 0
	j := 0
	for i < len(p) && j < len(path) {
		if p[i] == '*' {
			// abc.* abc.t, i = 4, j = 4
			if i == len(p)-1 {
				for ; j < len(path); j++ {
					if path[j] == filepath.Separator {
						return false
					}
				}
				return true
			}
			// ajax_*.html does not match ajax_/a.html
			if path[j] == filepath.Separator {
				return false
			}
			// abc*.html should match abc.ht.html
			s = i
			i++
			j++ // * should match at least one character
		} else if p[i] == '?' {
			if path[j] == filepath.Separator {
				return false
			}
			i++
			j++
		} else if p[i] == path[j] {
			i++
			j++
		} else if s != -1 && path[j] != filepath.Separator {
			i = s + 1
			j++
		} else {
			return false
		}
	}
	//log.Infof("len(p) %d i %d len(path) %d j %d", len(p), i, len(path), j)
	// both pattern and path reaches end
	return len(p) == i && len(path) == j
}
