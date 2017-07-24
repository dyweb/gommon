package util

import (
	"os"
	"strings"
	"testing"
	"io/ioutil"
)

func LoadDotEnv(t *testing.T) {
	b, err := ioutil.ReadFile(".env")
	if err != nil {
		t.Fatalf("failed to loade .env %v", err)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		kv := strings.SplitN(line, "=", 2)
		os.Setenv(kv[0], kv[1])
	}
}

// EnvAsMap returns environment variables as string map
// TODO: might cache it when package init, the problem of doing so is user might call os.Setenv, we also do this in test
func EnvAsMap() map[string]string {
	//https://coderwall.com/p/kjuyqw/get-environment-variables-as-a-map-in-golang
	envMap := make(map[string]string)

	for _, env := range os.Environ() {
		// NOTE: use SplitN instead of Split to handle the situation where value has =, i.e. key:FOO, value:bar1=bar2
		kv := strings.SplitN(env, "=", 2)
		envMap[kv[0]] = kv[1]
	}

	return envMap

}
