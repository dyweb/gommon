package generator

import (
	"bytes"
	"text/template"

	"github.com/dyweb/gommon/errors"
)

// NOTE: for now gommon only has logger

var _ GoConfig = (*LoggerConfig)(nil)

type LoggerConfig struct {
	Struct   string `yaml:"struct"`
	Receiver string `yaml:"receiver"`
	Field    string `yaml:"field"`
}

// RenderGommon returns nil, nil when there is nothing to render
// FIXME: gommon is an exception since it's took se
//func (c *ConfigFile) RenderGommon() ([]byte, error) {
//	// header
//	header := &bytes.Buffer{}
//	fmt.Fprintf(header, Header(Name, c.file))
//	fmt.Fprint(header, "\n")
//	if c.GoPackage != "" {
//		fmt.Fprintf(header, "package %s\n\n", c.GoPackage)
//	} else {
//		fmt.Fprintf(header, "package %s\n\n", c.pkg)
//	}
//	// body
//	body := &bytes.Buffer{}
//	// logger
//	if len(c.Loggers) > 0 {
//		fmt.Fprintln(header, "import dlog \"github.com/dyweb/gommon/log\"")
//		for _, l := range c.Loggers {
//			if err := l.RenderTo(body); err != nil {
//				return nil, err
//			}
//		}
//	}
//	if body.Len() == 0 {
//		return nil, nil
//	}
//	header.Write(body.Bytes())
//	//log.Debug(string(header.Bytes()))
//	// format go code
//	if formatted, err := format.Source(header.Bytes()); err != nil {
//		return formatted, errors.Wrap(err, "can't format generated code")
//	} else {
//		//log.Debugf("formatted len %d", len(formatted))
//		return formatted, nil
//	}
//}

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

func (c *LoggerConfig) Imports() []Import {
	return []Import{
		{"github.com/dyweb/gommon/log", "dlog"},
	}
}

func (c *LoggerConfig) FileName() string {
	return DefaultGeneratedFile
}

func (c *LoggerConfig) RenderBody(root string) ([]byte, error) {
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

func init() {
	structLoggerTmpl = template.Must(template.New(structLoggerTmplName).Parse(structLoggerTmplStr))
}
