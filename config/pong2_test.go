package config

import (
	"testing"

	"github.com/flosch/pongo2"
	asst "github.com/stretchr/testify/assert"
)

func TestRenderDocument(t *testing.T) {
	assert := asst.New(t)
	out, err := RenderDocumentString("{{ foo1 }} and {{ foo2 }} and {{ foo.a }}",
		pongo2.Context{"foo": map[string]interface{}{"a": 1}, "foo1": "bar", "foo2": 1})
	//t.Log(err)
	assert.Nil(err)
	assert.Equal("bar and 1 and 1", out)

}
