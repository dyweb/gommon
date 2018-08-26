// Package generator render go template, call external commands, generate gommon specific methods based on gommon.yml
package generator // import "github.com/dyweb/gommon/generator"

import (
	"github.com/dyweb/gommon/util/logutil"
)

const (
	GommonConfigFile = "gommon.yml"
	generatorName    = "gommon"
	GeneratedFile    = "gommon_generated.go"
)

var log = logutil.NewPackageLogger()
