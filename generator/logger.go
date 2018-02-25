package generator

import (
	"io"
	"text/template"

	"github.com/dyweb/gommon/errors"
)

type LoggerConfig struct {
	Struct   string `yaml:"struct"`
	Receiver string `yaml:"receiver"`
}

var structLoggerTmpl *template.Template

const structLoggerTmplName = "struct-logger"

const structLoggerTmplStr = `
func ({{.Receiver}} {{.Struct}}) SetLogger(logger *dlog.Logger) {
	{{.Receiver}}.log = logger
}

func ({{.Receiver}} {{.Struct}}) GetLogger() *dlog.Logger {
	return {{.Receiver}}.log
}

func ({{.Receiver}} {{.Struct}}) LoggerIdentity(justCallMe func() *dlog.Identity) *dlog.Identity {
	return justCallMe()
}
`

func (c *LoggerConfig) RenderTo(w io.Writer) error {
	if err := structLoggerTmpl.Execute(w, *c); err != nil {
		return errors.Wrap(err, "failed to render logger template")
	}
	return nil
}

func init() {
	structLoggerTmpl = template.Must(template.New(structLoggerTmplName).Parse(structLoggerTmplStr))
}
