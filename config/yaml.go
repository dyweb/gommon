package config

import (
	"bytes"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sync"
)

type YAMLConfig struct {
	data map[string]interface{}
	mu   sync.Mutex // TODO: may use RWMutex
}

type Context struct {
	vars map[string]interface{}
}

// SplitMultiDocument split a yaml file that contains multiple documents and
// (only) trim the first one if it is empty
func SplitMultiDocument(data []byte) [][]byte {
	docs := bytes.Split(data, []byte("---"))
	// check the first one, it could be empty
	if len(docs[0]) != 0 {
		return docs
	}
	return docs[1:]
}

func NewYAMLConfig() *YAMLConfig {
	c := new(YAMLConfig)
	c.data = make(map[string]interface{})
	return c
}

func (c *YAMLConfig) Parse(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := yaml.Unmarshal(data, &c.data)
	if err != nil {
		return errors.Wrap(err, "can't parse yaml to map[string]interface{}")
	}
	return nil
}

func (c *YAMLConfig) ParseMultiDocument() {
	// split the doc, parse by order, and add result to context so the following parser can use it
}

// func (c *YAMLConfig) ParseFile(path string) error {

// }
