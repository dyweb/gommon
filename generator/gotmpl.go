package generator

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/dyweb/gommon/util/genutil"
)

type GoTemplateConfig struct {
	Src  string      `yaml:"src"`
	Dst  string      `yaml:"dst"`
	Go   bool        `yaml:"go"`
	Data interface{} `yaml:"data"`
}

func (c *GoTemplateConfig) Render(root string) error {
	var (
		b   []byte
		buf bytes.Buffer
		err error
		t   *template.Template
	)
	//log.Infof("data is %v", c.Data)
	if b, err = ioutil.ReadFile(join(root, c.Src)); err != nil {
		return errors.Wrap(err, "can't read template file")
	}
	if t, err = template.New(c.Src).
		Funcs(template.FuncMap{
			"UcFirst": genutil.UcFirst,
		}).
		Parse(string(b)); err != nil {
		return errors.Wrap(err, "can't parse template")
	}
	buf.WriteString(genutil.DefaultHeader(join(root, c.Src)))
	if err = t.Execute(&buf, c.Data); err != nil {
		return errors.Wrap(err, "can't render template")
	}
	if c.Go {
		if b, err = genutil.Format(buf.Bytes()); err != nil {
			return errors.Wrap(err, "can't format as go code")
		}
	} else {
		b = buf.Bytes()
	}
	if err = fsutil.WriteFile(join(root, c.Dst), b); err != nil {
		return err
	}
	log.Debugf("rendered go tmpl %s to %s", join(root, c.Src), join(root, c.Dst))
	return nil
}
