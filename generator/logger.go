package generator

import (
	"io"
	"text/template"
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

func (c *LoggerConfig) RenderTo(w io.Writer) {
	structLoggerTmpl.Execute(w, *c)
}

func init() {
	structLoggerTmpl = template.Must(template.New(structLoggerTmplName).Parse(structLoggerTmplStr))
}
