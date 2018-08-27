package generator

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
)

var (
	_ Config   = (*GoTemplateConfig)(nil)
	_ GoConfig = (*GoTemplateConfig)(nil)
)

// TODO: go template does not need caller to add package header
type GoTemplateConfig struct {
	Src  string      `yaml:"src"`
	Dst  string      `yaml:"dst"`
	Go   bool        `yaml:"go"`
	Data interface{} `yaml:"data"`
}

func (c *GoTemplateConfig) IsGo() bool {
	return c.Go
}

// Imports is always empty, user should write imports in the template
func (c *GoTemplateConfig) Imports() []Import {
	return nil
}

func (c *GoTemplateConfig) FileName() string {
	return c.Dst
}

// TODO: go template does not need caller to add package header
func (c *GoTemplateConfig) RenderBody(root string) ([]byte, error) {
	if !c.IsGo() {
		return nil, errors.New("not go file, call Render instead")
	}
	// caller will do the formatting and add header
	return c.render(root)
}

func (c *GoTemplateConfig) Render(root string) error {
	if c.IsGo() {
		return errors.New("will generate go file, call RenderBody instead")
	}
	b, err := c.render(root)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString(DefaultHeader(join(root, c.Src)))
	buf.Write(b)
	if err = fsutil.WriteFile(join(root, c.Dst), buf.Bytes()); err != nil {
		return err
	}
	log.Debugf("rendered go tmpl %s to %s", join(root, c.Src), join(root, c.Dst))
	return nil
}

func (c *GoTemplateConfig) render(root string) ([]byte, error) {
	var (
		buf bytes.Buffer
		t   *template.Template
	)
	//log.Infof("data is %v", c.Data)
	b, err := ioutil.ReadFile(join(root, c.Src))
	if err != nil {
		return nil, errors.Wrap(err, "can't read template file")
	}
	if t, err = template.New(c.Src).Parse(string(b)); err != nil {
		return nil, errors.Wrap(err, "can't parse template")
	}
	if err = t.Execute(&buf, c.Data); err != nil {
		return nil, errors.Wrap(err, "can't render template")
	}
	return buf.Bytes(), nil
}
