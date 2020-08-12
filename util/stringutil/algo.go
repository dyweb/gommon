package stringutil

// algo.go implements common string algorithms like EditDistance.

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
