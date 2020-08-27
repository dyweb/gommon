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

// EnvToStruct decodes environment variable to struct
// TODO: it should works both normal struct and traceable struct config
// TODO: allow prefix, maybe pass config struct or as a method a config struct
func EnvToStruct(v interface{}) error {
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
		key := fieldNameToEnvKey(name)
		// TODO: use a env getter to avoid calling os.Getenv directly
		val := os.Getenv(key)
		switch fv.Kind() {
		case reflect.Int:
			if val == "" {
				fv.SetInt(0)
			} else {
				i, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("error parse key %s val %s for field %s: %w", key, val, name, err)
				}
				fv.SetInt(int64(i))
			}
		case reflect.Bool:
			if val == "" || val == "0" {
				fv.SetBool(false)
			} else {
				fv.SetBool(true)
			}
		case reflect.String:
			if val != "" {
				fv.SetString(val)
			}
		default:
			return fmt.Errorf("only int and bool field is supported got %s", fv.Kind())
		}
	}
	return nil
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
