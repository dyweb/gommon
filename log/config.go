package log

import (
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	Level           string                 `yaml:"level" json:"level"`
	Color           bool                   `yaml:"color" json:"color"`
	Source          bool                   `yaml:"source" json:"source"`
	ShowElapsedTime bool                   `yaml:"showElapsedTime" json:"showElapsedTime"`
	TimeFormat      string                 `yaml:"timeFormat" json:"timeFormat"`
	XXX             map[string]interface{} `yaml:",inline"`
}

func (c *Config) Validate() error {
	if c.XXX != nil {
		return errors.Errorf("undefined fields found %v", c.XXX)
	}
	if _, err := ParseLevel(c.Level, false); err != nil {
		return errors.Wrap(err, "invalid log level in config")
	}
	// valid time format
	t := time.Now()
	s := t.Format(c.TimeFormat)
	t2, err := time.Parse(c.TimeFormat, s)
	if err != nil || t.Unix() != t2.Unix() {
		return errors.Wrap(err, "invalid time format string in config")
	}
	return nil
}

func (log *Logger) ApplyConfig(c *Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if err := log.SetLevel(c.Level); err != nil {
		return err
	}
	if c.Source {
		log.EnableSourceLine()
	}
	log.Formatter.SetColor(c.Color)
	log.Formatter.SetElapsedTime(c.ShowElapsedTime)
	log.Formatter.SetTimeFormat(c.TimeFormat)
	// TODO: pkg filter should also be considered
	return nil
}
