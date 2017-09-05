package config

import (
	"github.com/dyweb/gommon/util"
	asst "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type logConfig struct {
	Level  string `yaml:"level"`
	Color  bool   `yaml:"color"`
	Source bool   `yaml:"source"`
}

type structuredConfig struct {
	Logging logConfig              `yaml:"logging"`
	Mode    string                 `yaml:"mode"`
	Base    string                 `yaml:"base"`
	Base2   string                 `yaml:"base2"`
	Base3   string                 `yaml:"base3"`
	XXX     map[string]interface{} `yaml:",inline"` // NOTE: this is used to catch unmatched fields
}

func TestYAMLConfig_ParseWithoutTemplate(t *testing.T) {
	assert := asst.New(t)
	var dat = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	c := NewYAMLConfig()
	err := c.ParseMultiDocument([]byte(dat))
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
	err = c.ParseMultiDocument([]byte(invalidDat))
	assert.NotNil(err)
	//log.Print(err.Error())
}

func TestLoadYAMLAsStruct(t *testing.T) {
	assert := asst.New(t)
	var conf structuredConfig
	assert.Nil(LoadYAMLAsStruct("testdata/structured.yml", &conf))
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

func TestYAMLConfig_ParseSingleDocument(t *testing.T) {
	cases := []struct {
		file string
	}{
		{"single_doc_no_vars"},
		{"single_doc_vars"},
	}
	c := NewYAMLConfig()
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			assert := asst.New(t)
			c.clear()
			doc := util.ReadFixture(t, "testdata/"+tc.file+".yml")
			assert.Nil(c.ParseSingleDocument(doc))
			// TODO: expect value, not just log
			t.Log(c.data)
		})
	}
}

func TestYAMLConfig_ParseMultiDocument(t *testing.T) {
	cases := []struct {
		file string
	}{
		{"multi_doc_single_vars"},
		{"multi_doc_multi_vars"},
	}
	c := NewYAMLConfig()
	//util.UseVerboseLog()
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			assert := asst.New(t)
			c.clear()
			doc := util.ReadFixture(t, "testdata/"+tc.file+".yml")
			assert.Nil(c.ParseMultiDocument(doc))
			t.Log(c.data)
		})
	}
	//util.DisableVerboseLog()
}

func TestYAMLConfig_Get(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocument(util.ReadFixture(t, "testdata/multi_doc_multi_vars.yml"))
	assert.Nil(err)
	//util.UseVerboseLog()
	assert.Equal("bar1", c.Get("vars.foo1"))
	assert.Equal(nil, c.Get("vars.that_does_not_exists"))
	// NOTE: top level keys other than vars are overwritten instead of merged
	assert.Equal(2, c.Get("foo"))
	//util.DisableVerboseLog()
}

func TestYAMLConfig_GetOrDefault(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocument(util.ReadFixture(t, "testdata/multi_doc_multi_vars.yml"))
	assert.Nil(err)
	assert.Equal("lalala", c.GetOrDefault("vars.oh_lala", "lalala"))
}

func TestYAMLConfig_GetOrFail(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocument(util.ReadFixture(t, "testdata/multi_doc_multi_vars.yml"))
	assert.Nil(err)
	_, err = c.GetOrFail("vars.oh_lala")
	assert.NotNil(err)
}

func TestYAMLConfig_Unmarshal(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocument(util.ReadFixture(t, "testdata/structured.yml"))
	assert.Nil(err)
	var conf structuredConfig
	// `vars` is always there even if it is not shown in config, sometimes we want user not to specify any fields we
	// didn't use, and use XXX (see above) to capture them, the removeVars flag would remove vars before Unmarshal and put
	// it back afterwards
	err = c.Unmarshal(&conf, false)
	assert.Nil(err)
	assert.NotNil(conf.XXX)
	assert.Equal("local", conf.Mode)
	// test envOr
	assert.Equal("base2", conf.Base2)
	assert.Equal(os.Getenv("HOME"), conf.Base3)
	err = c.Unmarshal(&conf, true)
	assert.Nil(err)
	assert.Nil(conf.XXX)
}

func TestYAMLConfig_UnmarshalKey(t *testing.T) {
	assert := asst.New(t)
	c := NewYAMLConfig()
	err := c.ParseMultiDocument(util.ReadFixture(t, "testdata/structured.yml"))
	assert.Nil(err)
	var conf logConfig
	err = c.UnmarshalKey("logging", &conf)
	assert.Nil(err)
	//t.Log(conf)
	assert.Equal("info", conf.Level)
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
