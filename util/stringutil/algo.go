package stringutil

// algo.go implements common string algorithms like EditDistance.

// CopySlice makes a deep copy of string slice.
// NOTE: it assumes the underlying string is immutable, i.e. the are not created using unsafe.
func CopySlice(src []string) []string {
	cp := make([]string, len(src))
	for i, s := range src {
		cp[i] = s
	}
	return cp
}

func EditDistance(word string, candidates []string, maxEdit int) {
	// TODO: return results group by distances [][]string should work when maxEdit is small
	// TODO: sort the response so it is stable?
	panic("not implemented")
}

// Shorter returns the shorter string or a if the length equals
func Shorter(a, b string) string {
	if len(a) <= len(b) {
		return a
	}
	return b
}

// Longer returns the longer string or b if the string equals
func Longer(a, b string) string {
	if len(a) >= len(b) {
		return a
	}
	return b
}
