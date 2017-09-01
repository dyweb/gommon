package util

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func ReadFixture(t *testing.T, path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read fixture %s: %v", path, err)
	}
	return b
}

func ReadJson(t *testing.T, path string, v interface{}) {
	b := ReadFixture(t, path)
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("can't unmarshal fixture %v", err)
	}
}
