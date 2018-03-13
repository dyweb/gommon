package testutil

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func PrintTidyJson(t *testing.T, data interface{}) {
	PrintTidyJsonTo(t, data, os.Stdout)
}

func PrintTidyJsonTo(t *testing.T, data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(data)
	if err != nil {
		t.Fatalf("print tidy json failed %v", err)
	}
}
