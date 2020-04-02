package generator_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/dyweb/gommon/util/genutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUcFirst(t *testing.T) {
	tmpl, err := template.New("test_ucfirst").Funcs(template.FuncMap{"UcFirst": genutil.UcFirst}).Parse(`{{ .foo | UcFirst }}`)
	require.Nil(t, err)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{
		"foo": "int64",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Int64", buf.String())
}
