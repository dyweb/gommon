package config

import (
	"github.com/dyweb/gommon/util"
	asst "github.com/stretchr/testify/assert"
	"testing"
)

var sampleMultiDoc = `
vars:
    influxdb_port: 8081
    databases:
        - influxdb
        - kairosdb
foo: 1
---
vars:
    kairosdb_port: 8080
    ssl: false
{% for db in vars.databases %}
    {{ db }}:
        name: {{ db }}
        ssl: {{ vars.ssl }}
{% endfor %}
foo: 2
`

func TestYAMLConfig_Parse(t *testing.T) {
	assert := asst.New(t)
	var dat = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	c := NewYAMLConfig()
	err := c.ParseMultiDocumentBytes([]byte(dat))
	assert.Nil(err)

	// NOTE： this is invalid yaml because when you use ` syntax to declare long string in Golang,
	// the indent are also included in the string, so this yaml has indent without any parent, which is invalid
	var invalidDat = `
	a: Easy!
	b:
	  c: 2
	  d: [3, 4]
	`
	// the print should show you the string has indent you may not be expecting
	//log.Print(invalidDat)
	err = c.ParseMultiDocumentBytes([]byte(invalidDat))
	assert.NotNil(err)
	//log.Print(err.Error())
}

func TestSplitMultiDocument(t *testing.T) {
	assert := asst.New(t)
	var multi = `---
time: 20:03:20
player: Sammy Sosa
action: strike (miss)
---
time: 20:03:47
player: Sammy Sosa
action: grand slam
`
	documents := SplitMultiDocument([]byte(multi))
	//for _, d := range  documents {
	//	t.Log(string(d[:]))
	//}
	assert.Equal(2, len(documents))
	documents = SplitMultiDocument([]byte("---"))
	assert.Equal(1, len(documents))
	// without the starting `---`
	var multi2 = `
time: 20:03:20
player: Sammy Sosa
action: strike (miss)
---
time: 20:03:47
player: Sammy Sosa
action: grand slam
`
	documents = SplitMultiDocument([]byte(multi2))
	assert.Equal(2, len(documents))
}

func TestYAMLConfig_ParseMultiDocumentBytes(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()

	// NOTE: use space instead of tab, YAML does not support tab
	// TODO: Add tab check in parser check, and tab inside quote should be allowed
	// WONTFIX: pongo2 render false to False, but Yaml spec support a lot of values http://yaml.org/type/bool.html
	var sampleUsePreviousVars = `
vars:
    influxdb_port: 8081
    databases:
        - influxdb
        - kairosdb
---
vars:
    kairosdb_port: 8080
    ssl: false
{% for db in vars.databases %}
    {{ db }}:
        name: {{ db }}
        ssl: {{ vars.ssl }}
{% endfor %}
`
	err := c.ParseMultiDocumentBytes([]byte(sampleUsePreviousVars))
	assert.Nil(err)
	// TODO: assert the value, need to use the dot syntax like Viper
	//assert.Equal(c.data)
	c.clear()

	var sampleUseCurrentVars = `
vars:
    foo: 1
bar:
    foo: {{ vars.foo }}
`
	err = c.ParseMultiDocumentBytes([]byte(sampleUseCurrentVars))
	assert.Nil(err)
	c.clear()

	// NOTE: I think HOME is set on most machines, at least travis?
	var sampleUseEnvironmentVars = `
vars:
    user: {{ envs.HOME }}
`
	err = c.ParseMultiDocumentBytes([]byte(sampleUseEnvironmentVars))
	assert.Nil(err)

}

func TestYAMLConfig_Get(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocumentBytes([]byte(sampleMultiDoc))
	assert.Nil(err)
	util.UseVerboseLog()
	assert.Equal(8081, c.Get("vars.influxdb_port"))
	assert.Equal(nil, c.Get("vars.that_does_not_exists"))
	// NOTE: top level keys other than vars are overwritten instead of merged
	assert.Equal(2, c.Get("foo"))
	util.DisableVerboseLog()
}

func TestYAMLConfig_GetOrFail(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocumentBytes([]byte(sampleMultiDoc))
	assert.Nil(err)
	_, err = c.GetOrFail("vars.oh_lala")
	assert.NotNil(err)
}

func TestYAMLConfig_GetOrDefault(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocumentBytes([]byte(sampleMultiDoc))
	assert.Nil(err)
	assert.Equal("lalala", c.GetOrDefault("vars.oh_lala", "lalala"))
}

func TestSearchMap(t *testing.T) {
	assert := asst.New(t)
	var m = make(map[string]interface{})
	var m2 = make(map[string]interface{})
	m["xephonk"] = m2
	m2["name"] = "xephonk"
	m2["port"] = 8080
	val, err := searchMap(m, []string{"xephonk", "name"})
	assert.Nil(err)
	assert.Equal("xephonk", val)
	_, err = searchMap(m, []string{"xephonk", "bar"})
	assert.NotNil(err)
}
