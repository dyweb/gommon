package config

import (
	"bytes"
	"sync"

	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type YAMLConfig struct {
	data   map[string]interface{}
	mu     sync.Mutex // TODO: may use RWMutex
	loader pongo2.TemplateLoader
	set    *pongo2.TemplateSet
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
	c.loader = pongo2.MustNewLocalFileSystemLoader("")
	c.set = pongo2.NewSet("gommon-yaml", c.loader)
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

// func (c *YAMLConfig) ParseFile(path string) error {

// }

func (c *YAMLConfig) ParseMultiDocumentBytes(data []byte) error {
	// split the doc, parse by order, and add result to context so the following parser can use it
	// docs := SplitMultiDocument(data)
	// for doc := range docs {
	// 	// rendered :=
	// }
	return nil
}

// TODO: this is quite duplicate with code in pongo2.go, but I think struct methods
// will be used more than those using default set? Or maybe the method should accept
// set as the first parameter
func (c *YAMLConfig) RenderDocumentString(tplStr string, context pongo2.Context) (string, error) {
	//pongo2.Context{} is just map[string]interface{}
	//FIXME: pongo2.FromString is not longer in the new API, must first create a set
	tpl, err := c.set.FromString(tplStr)
	if err != nil {
		return "", errors.Wrap(err, "can't parse template")
	}
	out, err := tpl.Execute(context)
	if err != nil {
		return "", errors.Wrap(err, "can'r render template")
	}
	return out, nil
}

func (c *YAMLConfig) RenderDocumentBytes(tplBytes []byte, context pongo2.Context) ([]byte, error) {
	tpl, err := c.set.FromBytes(tplBytes)
	var out []byte
	if err != nil {
		return out, errors.Wrap(err, "can't parse template")
	}
	out, err = tpl.ExecuteBytes(context)
	if err != nil {
		return out, errors.Wrap(err, "can'r render template")
	}
	return out, nil
}
