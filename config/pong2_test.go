package config

import (
	asst "github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderDocument(t *testing.T) {
	assert := asst.New(t)
	out, err := RenderDocument("{{ foo1 }} and {{ foo2 }}")
	//t.Log(err)
	assert.Nil(err)
	assert.Equal("bar and 1", out)

}
