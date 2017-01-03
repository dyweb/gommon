package config

import (
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	data map[string]interface{}
	mu   sync.Mutex // TODO: may use RWMutex
}

func New() *Config {
	c := new(Config)
	c.data = make(map[string]interface{})
	return c
}

func (c *Config) Parse(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := yaml.Unmarshal(data, &c.data)
	if err != nil {
		return errors.Wrap(err, "can't parse yaml to map[string]interface{}")
	}
	return nil
}

// func (c *Config) ParseFile(path string) error {

// }
