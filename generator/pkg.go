package generator

import (
	"fmt"
	"github.com/dyweb/gommon/util/logutil"
	"io"
)

var log = logutil.NewPackageLogger()

type Config struct {
	Loggers []LoggerConfig `yaml:"loggers"`
	pkg     string
}

func NewConfig(pkg string) *Config {
	return &Config{pkg: pkg}
}

// TODO: add error etc.
func (c *Config) Render(w io.Writer) {
	// TODO: should use go fmt, see Ayi gotmpl for example ...
	fmt.Fprintf(w, "package %s\n", c.pkg)
	fmt.Fprintln(w, "import dlog \"github.com/dyweb/gommon/log\"")
	for _, l := range c.Loggers {
		l.RenderTo(w)
	}
}
