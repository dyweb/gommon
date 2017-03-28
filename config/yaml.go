package config

import (
	"bytes"
	"sync"

	"fmt"
	"reflect"

	"github.com/dyweb/gommon/util"
	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// YAMLConfig is a thread safe struct for parse YAML file and get value
type YAMLConfig struct {
	vars   map[string]interface{}
	data   map[string]interface{}
	mu     sync.RWMutex // TODO: may use RWMutex
	loader pongo2.TemplateLoader
	set    *pongo2.TemplateSet
}

// SplitMultiDocument splits a yaml file that contains multiple documents and
// (only) trim the first one if it is empty
func SplitMultiDocument(data []byte) [][]byte {
	docs := bytes.Split(data, []byte("---"))
	// check the first one, it could be empty
	if len(docs[0]) != 0 {
		return docs
	}
	return docs[1:]
}

// NewYAMLConfig returns a config with internal map structure intialized
func NewYAMLConfig() *YAMLConfig {
	c := new(YAMLConfig)
	c.vars = make(map[string]interface{})
	c.data = make(map[string]interface{})
	c.loader = pongo2.MustNewLocalFileSystemLoader("")
	c.set = pongo2.NewSet("gommon-yaml", c.loader)
	return c
}

// clear is used by test for using one config object for several tets
func (c *YAMLConfig) clear() {
	c.vars = make(map[string]interface{})
	c.data = make(map[string]interface{})
}

func (c *YAMLConfig) ParseMultiDocumentBytes(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// split the doc, parse by order, add result to context so the following parser can use it
	docs := SplitMultiDocument(data)
	for _, doc := range docs {
		// TODO: pass environment variables
		pongoContext := pongo2.Context{
			"vars": c.vars,
			"env":  util.EnvAsMap(),
		}
		// TODO: explain the render twice logical for using vars in current document
		rendered, err := c.RenderDocumentBytes(doc, pongoContext)
		if err != nil {
			return errors.Wrap(err, "can't render template to yaml")
		}
		// TODO: need special flag/tag for this logging
		fmt.Printf("01-before\n%s", doc)
		fmt.Printf("01-after\n%s", rendered)
		tmpData := make(map[string]interface{})
		err = yaml.Unmarshal(rendered, &tmpData)
		if err != nil {
			return errors.Wrap(err, "can't parse rendered template yaml to map[string]interface{}")
		}
		// preserve the vars
		// TODO: move it to other function
		if varsInCurrentDocument, hasVars := tmpData["vars"]; hasVars {
			// NOTE: it's map[interface{}]interface{} instead of map[string]interface{}
			vars, ok := varsInCurrentDocument.(map[interface{}]interface{})
			if !ok {
				// TODO: test this
				// FIXME: why the log is map[string]interface{} while the cast is map[interface{}]interface{}
				return errors.Errorf("unable to cast %s to map[string]interface{}", reflect.TypeOf(varsInCurrentDocument))
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
		// NOTE: we don't need to assign c.vars to pongoContext again because it stores the reference to the map, not
		// the copy of the map
		rendered, err = c.RenderDocumentBytes(doc, pongoContext)
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
