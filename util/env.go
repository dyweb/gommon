package util

import (
	"os"
	"strings"
)

// EnvAsMap returns environment variables as string map
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
