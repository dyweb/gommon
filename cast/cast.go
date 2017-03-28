package cast

// ToStringMap converts a map to use string key, non string key will be ignored
func ToStringMap(interfaceMap map[interface{}]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(interfaceMap))

	for k, val := range interfaceMap {
		sk, ok := k.(string)
		if ok {
			m[sk] = val
		}
	}

	return m
}
