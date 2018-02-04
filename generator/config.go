package generator

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/pkg/errors"
)

type Config struct {
	Loggers     []LoggerConfig     `yaml:"loggers"`
	GoTemplates []GoTemplateConfig `yaml:"gotmpls"`
	// set when traversing the folders
	pkg  string
	file string
}

func NewConfig(pkg string, file string) *Config {
	return &Config{pkg: pkg, file: file}
}

func (c *Config) Render() ([]byte, error) {
	body := &bytes.Buffer{}
	header := &bytes.Buffer{}
	fmt.Fprintf(header, Header(generatorName, c.file))
	fmt.Fprintf(header, "package %s\n\n", c.pkg)
	if len(c.Loggers) > 0 {
		fmt.Fprintln(header, "import dlog \"github.com/dyweb/gommon/log\"")
		for _, l := range c.Loggers {
			l.RenderTo(body)
		}
	}
	// TODO: handle go template
	header.Write(body.Bytes())
	// format go code
	if formatted, err := format.Source(header.Bytes()); err != nil {
		return formatted, errors.Wrap(err, "can't format generated code")
	} else {
		return formatted, nil
	}
}
