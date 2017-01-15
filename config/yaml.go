package config

import (
	"bytes"
	"sync"

	"fmt"
	"reflect"

	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type YAMLConfig struct {
	vars   map[string]interface{}
	data   map[string]interface{}
	mu     sync.Mutex // TODO: may use RWMutex
	loader pongo2.TemplateLoader
	set    *pongo2.TemplateSet
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
	c.vars = make(map[string]interface{})
	c.data = make(map[string]interface{})
	c.loader = pongo2.MustNewLocalFileSystemLoader("")
	c.set = pongo2.NewSet("gommon-yaml", c.loader)
	return c
}

// clear is used by test in order not to create multiple config
func (c *YAMLConfig) clear() {
	c.vars = make(map[string]interface{})
	c.data = make(map[string]interface{})
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
	c.mu.Lock()
	defer c.mu.Unlock()
	// split the doc, parse by order, and add result to context so the following parser can use it
	docs := SplitMultiDocument(data)
	for _, doc := range docs {
		// TODO: pass environment variables
		rendered, err := c.RenderDocumentBytes(doc, pongo2.Context{"vars": c.vars})
		if err != nil {
			return errors.Wrap(err, "can't render template to yaml")
		}
		// TODO: user vars in current documents
		tmpData := make(map[string]interface{})
		fmt.Printf("01-before\n%s", doc)
		fmt.Printf("01-after\n%s", rendered)
		err = yaml.Unmarshal(rendered, &tmpData)
		if err != nil {
			return errors.Wrap(err, "can't parse rendered template yaml to map[string]interface{}")
		}
		// preserve the vars
		// TODO: rename to varsInCurrentDocument, hasVarsInCurrentDocument
		if varsRaw, ok := tmpData["vars"]; ok {
			// NOTE: it's map[interface{}]interface{} instead of map[string]interface{}
			vars, ok := varsRaw.(map[interface{}]interface{})
			if !ok {
				// TODO: test this
				return errors.Errorf("unable to cast %s to map[string]interface{}", reflect.TypeOf(varsRaw))
			}
			for k, v := range vars {
				// TODO: does YAML support non-string as key? if not, this assert is use less
				k, ok := k.(string)
				if !ok {
					// TODO: test this
					return errors.Errorf("unable to cast %s to string", reflect.TypeOf(k))
				}
				c.vars[k] = v
			}
		}
		// render again using vars in current document
		// TODO: if this document has no vars, then this render is not needed
		// TODO: use doc or previous render result
		rendered, err = c.RenderDocumentBytes(doc, pongo2.Context{"vars": c.vars})
		if err != nil {
			return errors.Wrap(err, "can't render template with vars in current document")
		}
		tmpData = make(map[string]interface{})
		fmt.Printf("02-before\n%s", doc)
		fmt.Printf("02-after\n%s", rendered)
		err = yaml.Unmarshal(rendered, &tmpData)
		if err != nil {
			// TODO: different message with previous error
			return errors.Wrap(err, "can't parse rendered template yaml to map[string]interface{}")
		}
		// put the data into c.data
		for k, v := range tmpData {
			c.data[k] = v
		}
	}
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
