// +build ignore

package structure

// Set is a map with string key and bool value
// It is not thread safe and modeled after https://github.com/deckarep/golang-set/blob/master/threadunsafe.go
type Set map[string]bool

// NewSet return a pointer of a set using arguments passed to the function
func NewSet(args ...string) *Set {
	// TODO: would this length for the map?
	m := make(Set, len(args))
	for _, v := range args {
		m[v] = true
	}
	return &m
}

// Cardinality return the size of the set
func (set *Set) Cardinality() int {
	return len(*set)
}

// Size is an alias for Cardinality
func (set *Set) Size() int {
	return len(*set)
}

// Contains check if a key is presented in the map, it does NOT check the bool value
func (set *Set) Contains(key string) bool {
	_, ok := (*set)[key]
	return ok
}

// Add add an element to set regardless of it is already in the set
func (set *Set) Add(key string) {
	(*set)[key] = true
}

// Equal check if two sets have exactly same elements
func (set *Set) Equal(other *Set) bool {
	if set.Cardinality() != other.Cardinality() {
		return false
	}
	for key := range *set {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}
