package testutil

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func ReadFixture(t *testing.T, path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read fixture %s: %v", path, err)
	}
	return b
}

func ReadJsonTo(t *testing.T, path string, v interface{}) {
	b := ReadFixture(t, path)
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("can't unmarshal fixture %s %v", path, err)
	}
}

func ReadYAMLTo(t *testing.T, path string, v interface{}) {
	b := ReadFixture(t, path)
	if err := yaml.Unmarshal(b, v); err != nil {
		t.Fatalf("can't unmarhsal YAML fixture %s %v", path, err)
	}
}

func WriteFixture(t *testing.T, path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0664)
	if err != nil {
		t.Fatalf("can't write fixture %s: %v", path, err)
	}
}
