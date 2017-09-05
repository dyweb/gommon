package util

import "sort"

// MergeStringMap merges the second map to first map, and keep the value of the second map if there is duplicate key
func MergeStringMap(dst map[string]interface{}, extra map[string]interface{}) {
	// TODO: rename extra to a more clear name
	for k, v := range extra {
		dst[k] = v
	}
}

// MapKeys returns keys of a map as string slice in an undermined order
// NOTE: you can not pass map[string]string because string is not interface{}, you need to convert it,
// but you can get the keys along the way of converting ...
func MapKeys(m map[string]interface{}) []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	// iterate map does not have order by golang spec
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapKeys returns keys of a map as string slice in an asc/desc order as specified
func MapSortedKeys(m map[string]interface{}, asc bool) []string {
	keys := MapKeys(m)
	if asc {
		sort.Strings(keys) // func Strings(a []string) { Sort(StringSlice(a)) }
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}
	return keys
}
