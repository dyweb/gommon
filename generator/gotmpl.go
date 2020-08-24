package generator

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/dyweb/gommon/util/stringutil"
	"golang.org/x/tools/imports"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/dyweb/gommon/util/genutil"
)

// GoTemplateConfig maps to gomtpls config in gommon.yml.
// It reads go template source from Src and writes to Dst using Data as data.
// The Go flag determines if the rendered template is formatted as go code.
type GoTemplateConfig struct {
	Src  string      `yaml:"src"`
	Dst  string      `yaml:"dst"`
	Go   bool        `yaml:"go"`
	Data interface{} `yaml:"data"`
}

func (c *GoTemplateConfig) Render(root string) error {
	var buf bytes.Buffer
	//log.Infof("data is %v", c.Data)
	src := join(root, c.Src)
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "can't read template file")
	}
	buf.WriteString(genutil.DefaultHeader(join(root, c.Src)))
	tmpl := GoCodeTemplate{
		Name:     src,
		Content:  string(b),
		Data:     c.Data,
		NoFormat: !c.Go,
		Funcs:    genutil.TemplateFuncMap(),
	}
	if err := RenderGoCodeTo(&buf, tmpl); err != nil {
		return err
	}
	dst := join(root, c.Dst)
	if err = fsutil.WriteFile(dst, buf.Bytes()); err != nil {
		return err
	}
	log.Debugf("rendered go tmpl %s to %s", src, dst)
	return nil
}

// ----------------------------------------------------------------------------
// GoCodeTemplate

type GoCodeTemplate struct {
	Name     string           // name used in error message e.g. generic-btree
	Content  string           // the actual template content
	Data     interface{}      // template data
	NoFormat bool             // disable calling goimports
	Funcs    template.FuncMap // additional template function map
}

func RenderGoCode(tmpl GoCodeTemplate) ([]byte, error) {
	// Sanity check
	if len(tmpl.Name) > len(tmpl.Content) {
		return nil, errors.Errorf("template name is longer than content, wrong order? shorter one is %s",
			stringutil.Shorter(tmpl.Name, tmpl.Content))
	}
	parsed, err := template.New(tmpl.Name).Parse(tmpl.Content)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, tmpl.Data); err != nil {
		return nil, errors.Wrapf(err, "error render template %s", tmpl.Name)
	}
	if tmpl.NoFormat {
		return buf.Bytes(), nil
	}
	formatted, err := FormatGo(buf.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "error format generated code")
	}
	return formatted, nil
}

func RenderGoCodeTo(dst io.Writer, tmpl GoCodeTemplate) error {
	// Sanity check
	if dst == nil {
		return errors.New("nil writer for RenderGoCode, forgot to check error when create file?")
	}
	b, err := RenderGoCode(tmpl)
	if err != nil {
		return err
	}
	_, err = dst.Write(b)
	return err
}

// ----------------------------------------------------------------------------
// Go Util

type GoStructDef struct {
	Name   string
	Fields []GoFieldDef
}

type GoFieldDef struct {
	Name string
	Type string
	Tag  string
}

// FormatGo formats go code using goimports without out fixing missing imports.
func FormatGo(src []byte) ([]byte, error) {
	opt := &imports.Options{
		Fragment:   false,
		AllErrors:  true,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: true,
	}
	return imports.Process("", src, opt)
}
