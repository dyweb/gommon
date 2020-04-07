package maputil

// interface.go provides helper related to map[string]interface{}

// MergeStringInterface creates a new map[string]interface{} that contains values from two maps.
// If there are duplicated keys, values from second map is preserved.
func MergeStringInterface(a, b map[string]interface{}) map[string]interface{} {
	c := make(map[string]interface{})
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}
