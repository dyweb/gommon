package cast

import (
	asst "github.com/stretchr/testify/assert"
	"testing"
)

type userConfig struct {
	Name string `yaml:"name" json:"name"`
}

func TestToStringMap(t *testing.T) {
	assert := asst.New(t)
	m := make(map[interface{}]interface{})
	m[1] = "123"
	m["milk"] = "cow"
	converted := ToStringMap(m)
	assert.Equal("cow", converted["milk"])
	assert.Equal(1, len(converted))
}

func TestStringMapToStructViaYAML(t *testing.T) {
	assert := asst.New(t)
	uc := &userConfig{}
	m := map[string]interface{}{"name": "jack"}
	assert.Nil(StringMapToStructViaYAML(m, uc))
	assert.Equal("jack", uc.Name)
}

func TestStringMapToStructViaJSON(t *testing.T) {
	assert := asst.New(t)
	uc := &userConfig{}
	m := map[string]interface{}{"name": "jack"}
	assert.Nil(StringMapToStructViaJSON(m, uc))
	assert.Equal("jack", uc.Name)
}
