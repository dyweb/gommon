package config

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/dyweb/gommon/cast"
	"github.com/dyweb/gommon/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// YAMLConfig is a thread safe struct for parse YAML file and get value
type YAMLConfig struct {
	vars         map[string]interface{}
	data         map[string]interface{}
	keyDelimiter string
	mu           sync.RWMutex
}

// SplitMultiDocument splits a yaml file that contains multiple documents and
// (only) trim the first one if it is empty
func SplitMultiDocument(data []byte) [][]byte {
	docs := bytes.Split(data, []byte(yamlDocumentSeparator))
	// check the first one, it could be empty
	if len(docs[0]) != 0 {
		return docs
	}
	return docs[1:]
}

// NewYAMLConfig returns a config with internal map structure initialized
func NewYAMLConfig() *YAMLConfig {
	c := new(YAMLConfig)
	c.clear()
	c.keyDelimiter = defaultKeyDelimiter
	return c
}

// clear is used by test for using one config object for several tests
// and is also used by constructor
func (c *YAMLConfig) clear() {
	c.vars = make(map[string]interface{})
	c.data = make(map[string]interface{})
}

func (c *YAMLConfig) Var(name string) interface{} {
	return c.vars[name]
}

func (c *YAMLConfig) Env(name string) string {
	return os.Getenv(name)
}

func (c *YAMLConfig) EnvOr(name string, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}

func (c *YAMLConfig) FuncMaps() template.FuncMap {
	return template.FuncMap{
		"env":   c.Env,
		"envOr": c.EnvOr,
		"var":   c.Var,
	}
}

func (c *YAMLConfig) ParseMultiDocument(data []byte) error {
	// split the doc, parse by order, add result to context so the following parser can use it
	docs := SplitMultiDocument(data)
	for _, doc := range docs {
		err := c.ParseSingleDocument(doc)
		if err != nil {
			return errors.Wrap(err, "error when parsing one of the documents")
		}
	}
	return nil
}

func (c *YAMLConfig) ParseSingleDocument(doc []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// we render the template twice, first time we use vars from previous documents and environment variables
	// second time, we use vars declared in this document, if any.

	// this is the first render
	rendered, err := c.renderDocument(doc)
	if err != nil {
		return errors.Wrap(err, "can't render template with previous documents' vars")
	}

	// TODO: change to trace level
	log.Debugf("01-before\n%s", doc)
	log.Debugf("01-after\n%s", rendered)

	tmpData := make(map[string]interface{})
	err = yaml.Unmarshal(rendered, &tmpData)
	if err != nil {
		return errors.Wrap(err, "can't parse rendered template yaml to map[string]interface{} after first render")
	}

	// preserve the vars
	if varsInCurrentDocument, hasVars := tmpData["vars"]; hasVars {
		// NOTE: it's map[interface{}]interface{} instead of map[string]interface{}
		// because how go-yaml handle the decoding
		varsI, ok := varsInCurrentDocument.(map[interface{}]interface{})
		if !ok {
			// TODO: test this, it seems if ok == false, then go-yaml should return already, which means ok should always be true?
			return errors.Errorf("unable to cast %s to map[interface{}]interface{}", reflect.TypeOf(varsInCurrentDocument))
		}
		vars := cast.ToStringMap(varsI)
		util.MergeStringMap(c.vars, vars)

		// render again using vars in current document
		rendered, err = c.renderDocument(doc)
		if err != nil {
			return errors.Wrap(err, "can't render template with vars in current document")
		}

		log.Debugf("02-before\n%s", doc)
		log.Debugf("02-after\n%s", rendered)

		tmpData = make(map[string]interface{})
		err = yaml.Unmarshal(rendered, &tmpData)
		if err != nil {
			return errors.Wrap(err, "can't parse rendered template yaml to map[string]interface{} after second render")
		}
	}

	// put the data into c.data
	for k, v := range tmpData {
		c.data[k] = v
	}
	// NOTE: vars are merged instead of overwritten like other top level keys
	c.data["vars"] = c.vars
	return nil
}

func (c *YAMLConfig) Get(key string) interface{} {
	val, err := c.GetOrFail(key)
	if err != nil {
		// TODO: so Get does not do any error handling? what did viper do?
		log.Debugf("can't get key %s due to %s", key, err.Error())
	}
	return val
}

func (c *YAMLConfig) GetOrDefault(key string, defaultVal interface{}) interface{} {
	val, err := c.GetOrFail(key)
	if err != nil {
		log.Debugf("use default %v for key %s due to %s", defaultVal, key, err.Error())
		return defaultVal
	}
	return val
}

func (c *YAMLConfig) GetOrFail(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	path := strings.Split(key, c.keyDelimiter)
	val, err := searchMap(c.data, path)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *YAMLConfig) Unmarshal(structuredConfig interface{}, removeVars bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: maybe we no longer need to keep vars in data?
	if removeVars {
		delete(c.data, "vars")
		defer func() { c.data["vars"] = c.vars }()
	}
	out, err := yaml.Marshal(c.data)
	// FIXME: those Wrapf should be Wrap
	if err != nil {
		return errors.Wrapf(err, "can't marshal data back to yaml")
	}
	err = yaml.Unmarshal(out, structuredConfig)
	if err != nil {
		return errors.Wrapf(err, "can't unmarshal data to specified config struct")
	}
	return nil
}

func (c *YAMLConfig) UnmarshalKey(key string, structuredConfig interface{}) error {
	val, err := c.GetOrFail(key)
	if err != nil {
		return errors.Wrapf(err, "can't get the val for unmarshal")
	}
	out, err := yaml.Marshal(val)
	if err != nil {
		return errors.Wrapf(err, "can't marshal data back to yaml")
	}
	err = yaml.Unmarshal(out, structuredConfig)
	if err != nil {
		return errors.Wrapf(err, "can't unmarshal data to specified config struct")
	}
	return nil
}

// TODO: what if the user use key that includes `.`, like {"github.com": 123}
func searchMap(src map[string]interface{}, path []string) (interface{}, error) {
	var result interface{}
	if len(path) == 0 {
		return result, errors.New("path is empty, at least provide one segment")
	}
	result = src
	previousPath := ""
	for i := 0; i < len(path); i++ {
		key := path[i]
		//log.Debug(key)
		//log.Debug(result)
		if reflect.TypeOf(result).Kind() != reflect.Map {
			return nil, errors.Errorf("%s is not a map but %s, %v", previousPath, reflect.TypeOf(result), result)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			m = cast.ToStringMap(result.(map[interface{}]interface{}))
		}
		// FIXED: this is a tricky problem, if you use `:` here, you create a new local variable instead update the one
		// outside the loop, that's all the Debug(result) for
		//result, ok := m[key]
		result, ok = m[key]
		if !ok {
			return result, errors.Errorf("key: %s does not exists in path: %s, val: %v", key, previousPath, m)
		}
		//log.Debug(result)
		previousPath += key
	}
	return result, nil
}

func (c *YAMLConfig) renderDocument(tplBytes []byte) ([]byte, error) {
	tmpl, err := template.New(defaultTemplateName).Funcs(c.FuncMaps()).Parse(string(tplBytes[:]))
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse template")
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, "")
	if err != nil {
		return nil, errors.Wrapf(err, "can't render template")
	}
	return b.Bytes(), nil
}
