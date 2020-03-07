package testutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// env.go load .env file during testing

// LoadDotEnv load .env in current directory into environment variable, line start with # are comments
// It is modeled after https://github.com/motdotla/dotenv
func LoadDotEnv(t *testing.T) {
	LoadDotEnvFromFile(t, ".env")
}

func LoadDotEnvFromFile(t *testing.T, path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load .env %v", err)
	}
	LoadDotEnvFrom(t, bytes.NewReader(b))
}

func LoadDotEnvFrom(t *testing.T, r io.Reader) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read .env content %v", err)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) == 1 {
			if err := os.Setenv(kv[0], ""); err != nil {
				t.Fatalf("failed to set env %s %v", kv[0], err)
			}
		}
		if len(kv) == 2 {
			if err := os.Setenv(kv[0], kv[1]); err != nil {
				t.Fatalf("failed to set env %s=%s %v", kv[0], kv[1], err)
			}
		}
	}
}
