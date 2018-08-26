package generator

import (
	"io"
	"text/template"

	"github.com/dyweb/gommon/errors"
)

// NOTE: for now gommon only has logger

type LoggerConfig struct {
	Struct   string `yaml:"struct"`
	Receiver string `yaml:"receiver"`
	Field    string `yaml:"field"`
}

var structLoggerTmpl *template.Template

const structLoggerTmplName = "struct-logger"

const structLoggerTmplStr = `
func ({{.Receiver}} {{.Struct}}) SetLogger(logger *dlog.Logger) {
	{{.Receiver}}.{{.Field}} = logger
}

func ({{.Receiver}} {{.Struct}}) GetLogger() *dlog.Logger {
	return {{.Receiver}}.{{.Field}}
}

func ({{.Receiver}} {{.Struct}}) LoggerIdentity(justCallMe func() *dlog.Identity) *dlog.Identity {
	return justCallMe()
}
`

func (c *LoggerConfig) RenderTo(w io.Writer) error {
	// NOTE: (at15) for backward compatibility, will remove it once refactor on generator is done
	if c.Field == "" {
		c.Field = "log"
	}
	if err := structLoggerTmpl.Execute(w, *c); err != nil {
		return errors.Wrap(err, "failed to render logger template")
	}
	return nil
}

func init() {
	structLoggerTmpl = template.Must(template.New(structLoggerTmplName).Parse(structLoggerTmplStr))
}
