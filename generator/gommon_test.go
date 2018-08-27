package generator

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
	requir "github.com/stretchr/testify/require"
)

func TestLoggerConfig_RenderTo(t *testing.T) {
	cfgs := []struct {
		name     string
		c        LoggerConfig
		rendered string
	}{
		{"default field", LoggerConfig{"*YAMLConfig", "c", ""}, `
func (c *YAMLConfig) SetLogger(logger *dlog.Logger) {
	c.log = logger
}

func (c *YAMLConfig) GetLogger() *dlog.Logger {
	return c.log
}

func (c *YAMLConfig) LoggerIdentity(justCallMe func() *dlog.Identity) *dlog.Identity {
	return justCallMe()
}
`},
		{"specified field", LoggerConfig{"*YAMLConfig", "c", "logger"}, `
func (c *YAMLConfig) SetLogger(logger *dlog.Logger) {
	c.logger = logger
}

func (c *YAMLConfig) GetLogger() *dlog.Logger {
	return c.logger
}

func (c *YAMLConfig) LoggerIdentity(justCallMe func() *dlog.Identity) *dlog.Identity {
	return justCallMe()
}
`},
	}
	for _, cfg := range cfgs {
		t.Run(cfg.name, func(t *testing.T) {
			assert := asst.New(t)
			require := requir.New(t)

			b, err := cfg.c.RenderBody("")
			require.Nil(err)
			assert.Equal(cfg.rendered, string(b))
		})
	}
}
