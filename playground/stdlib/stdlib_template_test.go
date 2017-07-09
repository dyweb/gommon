package stdlib

import (
	"os"
	"text/template"
	"testing"
	"bytes"

	asst "github.com/stretchr/testify/assert"
)

func TestTemplate_Funcs(t *testing.T) {
	assert := asst.New(t)
	const tmplText = `
	{{env "HOME"}}
	{{var "foo"}}
`
	vars := map[string]string{"foo": "bar"}

	funcMap := template.FuncMap{
		"env": func(name string) string {
			return os.Getenv(name)
		},
		"var": func(name string) string {
			return vars[name]
		},
	}

	tmpl, err := template.New("funcs test").Funcs(funcMap).Parse(tmplText)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	assert.Nil(tmpl.Execute(&b, ""))
	rendered := b.String()
	assert.Contains(rendered, os.Getenv("HOME"), "env function in template should be called")
	assert.Contains(rendered, "bar")

	// now we update vars
	vars["foo"] = "bar2"
	b.Reset()
	assert.Nil(tmpl.Execute(&b, ""))
	rendered = b.String()
	assert.Contains(rendered, "bar2", "var function should be using the updated value")
}

func TestTemplate_Range(t *testing.T) {
	assert := asst.New(t)
	// NOTE: the `-` is used to remove the following new line https://golang.org/pkg/text/template/#hdr-Text_and_spaces
	// FIXME: - will also remove space which would break yaml indent
	const tmplText = `
{{ range $name := var "databases" -}}
{{ $name -}}:
{{ $db := var $name }}
    name: {{ $name }}
    type: {{ $db.type -}}
{{ end }}
`

	vars := map[string]interface{}{"foo": "barr"}
	vars["databases"] = []string{"cassandra", "mysql", "xephonk"}
	vars["cassandra"] = map[string]string{"type": "nosql"}
	vars["mysql"] = map[string]string{"type": "sql"}
	vars["xephonk"] = map[string]string{"type": "tsdb"}

	funcMap := template.FuncMap{
		"env": func(name string) string {
			return os.Getenv(name)
		},
		"var": func(name string) interface{} {
			return vars[name]
		},
	}

	const rendered = `
cassandra:

    name: cassandra
    type: nosqlmysql:

    name: mysql
    type: sqlxephonk:

    name: xephonk
    type: tsdb
`
	tmpl, err := template.New("range test").Funcs(funcMap).Parse(tmplText)
	assert.Nil(err)
	var b bytes.Buffer
	err = tmpl.Execute(&b, "")
	//t.Log(err)
	assert.Nil(err)
	assert.Equal(rendered, b.String())
}
