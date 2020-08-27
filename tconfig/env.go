package tconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// env.go defines environment variable related operations

// EnvReader reads environment variable by key and returns empty string if it does not exists
type EnvReader func(key string) (value string)

// EvalBool converts a string value to boolean value.
type EvalBool func(s string) bool

type EnvToStructConfig struct {
	Prefix string
	Env    EnvReader // Env determines how to get environment variable, os.Getenv is used when set to nil
	Bool   EvalBool  // Bool determines how to evaluate a string to bool, DefaultEvalBool is used when set to nil
}

// a private default so other packages can't change it
var defaultEnvToStructConfig = DefaultEnvToStructConfig()

func DefaultEnvToStructConfig() EnvToStructConfig {
	return EnvToStructConfig{
		Env:  os.Getenv,
		Bool: DefaultEvalBool(),
	}
}

// a private default so other packages can't change it
var defaultEvalBool = DefaultEvalBool()

// DefaultEvalBool treats empty string, 0, false, FALSE as false.
func DefaultEvalBool() EvalBool {
	return func(s string) bool {
		if s == "" || s == "0" || s == "false" || s == "FALSE" {
			return false
		}
		return true
	}
}

// TODO: it should works both normal struct and traceable struct config
func (c *EnvToStructConfig) To(v interface{}) error {
	envReader := c.Env
	if envReader == nil {
		envReader = os.Getenv
	}
	evalBool := c.Bool
	if evalBool == nil {
		evalBool = defaultEvalBool
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v must be a non nil pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("can only decode to a struct got %s", rv.Kind())
	}
	rt := rv.Type()
	nFields := rv.NumField()
	for i := 0; i < nFields; i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		name := ft.Name
		key := c.Prefix + fieldNameToEnvKey(name)
		val := envReader(key)
		switch fv.Kind() {
		case reflect.Int:
			if val != "" {
				i, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("invalid int key %s val %s for field %s: %w", key, val, name, err)
				}
				fv.SetInt(int64(i))
			}
		case reflect.Bool:
			fv.SetBool(evalBool(val))
		case reflect.String:
			if val != "" {
				fv.SetString(val)
			}
		default:
			return fmt.Errorf("only int, bool and string fields are supported got %s", fv.Kind())
		}
	}
	return nil
}

// EnvToStruct decodes environment variable to struct using DefaultEnvToStructConfig
func EnvToStruct(v interface{}) error {
	return defaultEnvToStructConfig.To(v)
}

func fieldNameToEnvKey(name string) string {
	var b strings.Builder
	for i, r := range name {
		if unicode.IsUpper(r) && i != 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// ----------------------------------------------------------------------------
// Env Util, they are not that related to tconfig but save some typing in adhoc code

// EnvInt returns decoded int or 0 if the key does not exists or contains invalid value.
func EnvInt(key string) int {
	return EnvIntDefault(key, 0)
}

// EnvIntDefault returns decoded int or defaultValue if the key does not exists or contains invalid value.
func EnvIntDefault(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}
