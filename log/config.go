package log

import (
	"github.com/pkg/errors"
)

// TODO: timestamp format, time duration format
type Config struct {
	Level  string                 `yaml:"level" json:"level"`
	Color  bool                   `yaml:"color" json:"color"`
	Source bool                   `yaml:"source" json:"source"`
	XXX    map[string]interface{} `yaml:",inline"`
}

func (c *Config) Validate() error {
	if c.XXX != nil {
		return errors.Errorf("undefined fields %v", c.XXX)
	}
	if _, err := ParseLevel(c.Level, false); err != nil {
		return errors.Wrap(err, "invalid log level in config")
	}
	return nil
}
