package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func TestConfig_Validate(t *testing.T) {
	assert := asst.New(t)
	b := readFixture(t, "config.example.yml")
	var c Config
	assert.Nil(yaml.Unmarshal(b, &c))
	assert.Nil(c.Validate())
	l := NewLogger()
	assert.Nil(l.ApplyConfig(&c))
	e := l.NewEntryWithPkg("pkg1")
	e.Info("test config")
}

// FIXME: this is copied from util library to avoid import cycle
func readFixture(t *testing.T, path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read fixture %s: %v", path, err)
	}
	return b
}
