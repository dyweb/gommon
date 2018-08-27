// Package envutil wraps environment variable related operations
package envutil

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/maputil"
)

func LoadDotEnv(file string) error {
	if file == "" {
		file = ".env"
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "error load .env from %s", file)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) == 1 {
			os.Setenv(kv[0], "")
		}
		if len(kv) == 2 {
			os.Setenv(kv[0], kv[1])
		}
	}
	return nil
}

// EnvMap returns environment variables as string map without cache
func EnvMap() map[string]string {
	//https://coderwall.com/p/kjuyqw/get-environment-variables-as-a-map-in-golang
	envMap := make(map[string]string)

	for _, env := range os.Environ() {
		// NOTE: use SplitN instead of Split to handle the situation where value has =, i.e. key:FOO, value:bar1=bar2
		kv := strings.SplitN(env, "=", 2)
		envMap[kv[0]] = kv[1]
	}
	return envMap
}

var initEnvMap map[string]string

// InitEnvMap returns the environment map when package is initialized, it may not be accurate due to other package
// modifying environment in init or the package is initialized via Go's plugin
func InitEnvMap() map[string]string {
	return maputil.CopyStringMap(initEnvMap)
}

func init() {
	initEnvMap = EnvMap()
}
