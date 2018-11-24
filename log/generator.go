package log

import (
	"bytes"
	"text/template"

	"github.com/dyweb/gommon/errors"
)

var structLoggerTmpl *template.Template

// StructLoggerConfig is used to generate methods on struct for get identity using runtime,
// it also generates getter and setter
type StructLoggerConfig struct {
	Struct   string `yaml:"struct"`
	Receiver string `yaml:"receiver"`
	Field    string `yaml:"field"`
}

const structLoggerTmplStr = `
func ({{.Receiver}} {{.Struct}}) SetLogger(logger *dlog.Logger) {
	{{.Receiver}}.{{.Field}} = logger
}

func ({{.Receiver}} {{.Struct}}) GetLogger() *dlog.Logger {
	return {{.Receiver}}.{{.Field}}
}

func ({{.Receiver}} {{.Struct}}) LoggerIdentity(justCallMe func() dlog.Identity) dlog.Identity {
	return justCallMe()
}
`

const structLoggerTmplName = "struct-logger"

func (c *StructLoggerConfig) Render() ([]byte, error) {
	if structLoggerTmpl == nil {
		tpl, err := template.New(structLoggerTmplName).Parse(structLoggerTmplStr)
		if err != nil {
			return nil, errors.Wrap(err, "error parse template")
		}
		structLoggerTmpl = tpl
	}
	// NOTE: (at15) for backward compatibility, will remove it once refactor on generator is done
	if c.Field == "" {
		c.Field = "log"
	}
	var buf bytes.Buffer
	if err := structLoggerTmpl.Execute(&buf, *c); err != nil {
		return nil, errors.Wrap(err, "failed to render logger template")
	}
	return buf.Bytes(), nil
}
