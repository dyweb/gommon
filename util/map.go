package util

// MergeStringMap merges the second map to first map, and keep the value of the second map if there is duplicate key
func MergeStringMap(dst map[string]interface{}, extra map[string]interface{}) {
	// TODO: rename extra to a more clear name
	for k, v := range extra {
		dst[k] = v
	}
}
