package generator

import (
	"fmt"
	"go/format"

	"bytes"
	"github.com/dyweb/gommon/util/logutil"
	"github.com/pkg/errors"
)

const generatorName = "gommon"

var log = logutil.NewPackageLogger()

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
	b := &bytes.Buffer{}
	fmt.Fprintf(b, Header(generatorName, c.file))
	fmt.Fprintf(b, "package %s\n\n", c.pkg)
	fmt.Fprintln(b, "import dlog \"github.com/dyweb/gommon/log\"")
	for _, l := range c.Loggers {
		l.RenderTo(b)
	}
	// format go code
	if formatted, err := format.Source(b.Bytes()); err != nil {
		return formatted, errors.Wrap(err, "can't format generated code")
	} else {
		return formatted, nil
	}
}
