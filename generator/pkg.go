// Package generator render go template, call external commands, generate gommon specific methods based on gommon.yml
package generator

import (
	"github.com/dyweb/gommon/util/logutil"
)

const (
	configFile    = "gommon.yml"
	generatorName = "gommon"
	generatedFile = "gommon_generated.go"
)

var log = logutil.NewPackageLogger()
