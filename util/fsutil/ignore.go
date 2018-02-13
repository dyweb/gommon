package fsutil

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
type WildcardPattern string

// TODO: it is not working ....
// TODO: test on windows
func (p WildcardPattern) ShouldIgnore(path string) bool {
	// the pattern would always be no greater than path due to our limited features
	if len(p) > len(path) {
		return false
	}
	i := 0
	j := 0
	for i < len(p) && j < len(path) {
		if p[i] == '*' {
			// abc.*
			if i == len(p)-1 {
				// TODO: len or len - 1?
				return len(path) > j
			}
			// *.html
			i++

		}
		if p[i] == '?' || p[i] == path[j] {
			i++
			j++
		}
	}
	// both pattern and path reaches end
	return len(p)-1 == i && len(path)-1 == j
}
