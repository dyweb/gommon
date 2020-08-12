package genutil

import (
	htmltemplate "html/template"
	texttemplate "text/template"
	"unicode"

	"github.com/dyweb/gommon/util/stringutil"
)

// template_funcs.go defines common template functions used in go template

// TemplateFuncMap returns a new func map that includes all template func in this page.
func TemplateFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"UcFirst": stringutil.UcFirst,
		"LcFirst": stringutil.LcFirst,
	}
}

// TextTemplateFuncMap returns func map for text/template
func TextTemplateFuncMap() texttemplate.FuncMap {
	return TemplateFuncMap()
}

// HTMLTemplateFuncMap returns func map for html/template
// TODO: maybe we should have some extra html specific helpers
func HTMLTemplateFuncMap() htmltemplate.FuncMap {
	return TemplateFuncMap()
}

// UcFirst changes first character to upper case.
// Deprecated: use stringuitl.UcFirst
func UcFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// SnakeToCamel converts snake_case to CamelCase.
// Deprecated: use stringutil.SnakeToCamel
func SnakeToCamel(s string) string {
	src := []rune(s)
	var dst []rune
	toUpper := true
	for _, r := range src {
		if r == '_' {
			toUpper = true
			continue
		}

		r2 := r
		if toUpper {
			r2 = unicode.ToUpper(r)
			toUpper = false
		}
		dst = append(dst, r2)
	}
	return string(dst)
}

// LcFirst changes first character to lower case.
// Deprecated: use stringutil.SnakeToCamel
func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
