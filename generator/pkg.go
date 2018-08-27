// Package generator render go template, call external commands, generate gommon specific methods based on gommon.yml
package generator // import "github.com/dyweb/gommon/generator"

import (
	"github.com/dyweb/gommon/util/logutil"
)

const (
	// PkgName TODO: is reserved for the incoming library package (liblibrary?) to allow control logger etc.
	PkgName              = "gommon"
	GommonConfigFile     = "gommon.yml"
	Name                 = "gommon"
	GeneratedFile        = "gommon_generated.go"
	DefaultGeneratedFile = "gommon_generated.go"
)

var log = logutil.NewPackageLogger()

type Config interface {
	IsGo() bool
	Render(root string) error
}

type Import struct {
	Pkg   string
	Alias string
}

type GoConfig interface {
	// Imports that will be put at top of file
	Imports() []Import
	// FileName is the name the generator wants the caller to use when saving content
	FileName() string
	// RenderBody returns the body without imports
	RenderBody(root string) ([]byte, error)
}

type ConfigFile struct {
	// Loggers is helper methods on struct for gommon/log to build a tree for logger, this is subject to change
	Loggers []LoggerConfig `yaml:"loggers"`
	// GoTemplates is templates written in go's text/template format, they are mainly used to generate go source file
	GoTemplates []GoTemplateConfig `yaml:"gotmpls"`
	// Noodles is the config for embedding assets by generating go file with a large byte slice
	Noodles []NoodleConfig `yaml:"noodles"`
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
