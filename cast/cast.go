package cast

import (
	"encoding/json"

	"github.com/dyweb/gommon/errors"
	"gopkg.in/yaml.v2"
)

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

func StringMapToStructViaYAML(m map[string]interface{}, s interface{}) error {
	b, err := yaml.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "can't marshal to YAML")
	}
	if err = yaml.Unmarshal(b, s); err != nil {
		return errors.Wrap(err, "can't unmarshal YAML to struct")
	}
	return nil
}

func StringMapToStructViaJSON(m map[string]interface{}, s interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "can't marshal to JSON")
	}
	if err = json.Unmarshal(b, s); err != nil {
		return errors.Wrap(err, "can't unmarshal JSON to struct")
	}
	return nil
}
