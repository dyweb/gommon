package generator

import (
	"bytes"
	"github.com/pkg/errors"
	"go/format"
	"io/ioutil"
	"text/template"
)

type GoTemplateConfig struct {
	Src  string      `yaml:"src"`
	Dst  string      `yaml:"dst"`
	Data interface{} `yaml:"data"`
}

func (c *GoTemplateConfig) Render(root string) error {
	var (
		b   []byte
		buf bytes.Buffer
		err error
		t   *template.Template
	)
	if b, err = ioutil.ReadFile(join(root, c.Src)); err != nil {
		return errors.Wrap(err, "can't read template file")
	}
	if t, err = template.New("main").Parse(string(b)); err != nil {
		return errors.Wrap(err, "can't parse template")
	}
	buf.WriteString(Header(generatorName, join(root, c.Src)))
	if err = t.Execute(&buf, c.Data); err != nil {
		return errors.Wrap(err, "can't render template")
	}
	// TODO: support non go code using go template a.k.a no go format
	if b, err = format.Source(buf.Bytes()); err != nil {
		return errors.Wrap(err, "can't format as go code")
	}
	if err = WriteFile(join(root, c.Dst), b); err != nil {
		return err
	}
	log.Debugf("rendered go tmpl %s to %s", join(root, c.Src), join(root, c.Dst))
	return nil
}
