package generator

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/dyweb/gommon/errors"
)

type Config interface {
	IsGo() bool
	Render(root string) error
}

type GoConfig interface {
	// Imports that will be put at top of file
	Imports() []string
	// FileName is the name the generator wants the caller to use when saving content
	FileName() string
	// RenderBody returns the body without imports
	RenderBody(root string) ([]byte, error)
}

type ConfigFile struct {
	// Loggers is helper methods on struct for gommon/log to build a tree for logger, this is subject to change
	Loggers     []LoggerConfig     `yaml:"loggers"`
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

// RenderGommon returns nil, nil when there is nothing to render
func (c *ConfigFile) RenderGommon() ([]byte, error) {
	// header
	header := &bytes.Buffer{}
	fmt.Fprintf(header, Header(generatorName, c.file))
	fmt.Fprint(header, "\n")
	if c.GoPackage != "" {
		fmt.Fprintf(header, "package %s\n\n", c.GoPackage)
	} else {
		fmt.Fprintf(header, "package %s\n\n", c.pkg)
	}
	// body
	body := &bytes.Buffer{}
	// logger
	if len(c.Loggers) > 0 {
		fmt.Fprintln(header, "import dlog \"github.com/dyweb/gommon/log\"")
		for _, l := range c.Loggers {
			if err := l.RenderTo(body); err != nil {
				return nil, err
			}
		}
	}
	if body.Len() == 0 {
		return nil, nil
	}
	header.Write(body.Bytes())
	//log.Debug(string(header.Bytes()))
	// format go code
	if formatted, err := format.Source(header.Bytes()); err != nil {
		return formatted, errors.Wrap(err, "can't format generated code")
	} else {
		//log.Debugf("formatted len %d", len(formatted))
		return formatted, nil
	}
}
