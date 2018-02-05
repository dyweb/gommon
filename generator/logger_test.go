package generator

import (
	"bytes"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestLoggerConfig_RenderTo(t *testing.T) {
	assert := asst.New(t)
	c := LoggerConfig{Struct: "*YAMLConfig", Receiver: "c"}
	var rendered = `
func (c *YAMLConfig) SetLogger(logger *dlog.Logger) {
	c.log = logger
}

func (c *YAMLConfig) GetLogger() *dlog.Logger {
	return c.log
}

func (c *YAMLConfig) LoggerIdentity(justCallMe func() *dlog.Identity) *dlog.Identity {
	return justCallMe()
}
`
	var b bytes.Buffer
	c.RenderTo(&b)
	assert.Equal(rendered, string(b.Bytes()))
}
