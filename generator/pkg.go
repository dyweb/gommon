// Package generator render go template, call external commands, generate gommon specific methods based on gommon.yml
package generator // import "github.com/dyweb/gommon/generator"

import (
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/noodle"
)

const (
	GommonIgnoreFile     = ".gommonignore"
	GommonConfigFile     = "gommon.yml"
	DefaultGeneratedFile = "gommon_generated.go"
)

var (
	logReg = dlog.NewRegistry()
	log    = logReg.Logger()
)

type ConfigFile struct {
	// Loggers is helper methods on struct for gommon/log to build a tree for logger, this is subject to change
	Loggers []dlog.StructLoggerConfig `yaml:"loggers"`
	// GoTemplates is templates written in go's text/template format, they are mainly used to generate go source file
	GoTemplates []GoTemplateConfig `yaml:"gotmpls"`
	// Noodles is the config for embedding assets by generating go file with a large byte slice
	Noodles []noodle.EmbedConfig `yaml:"noodles"`
	// Shells is shell commands to be executed
	Shells []ShellConfig `yaml:"shells"`
	// GoPackage override folder name for generated go file
	GoPackage string `yaml:"go_package"`

	// set when traversing the folders
	pkg  string
	file string
}

func NewConfigFile(pkg string, file string) *ConfigFile {
	return &ConfigFile{pkg: pkg, file: file}
}
