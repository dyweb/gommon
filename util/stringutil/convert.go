package stringutil

import "unicode"

// convert.go converts string and strings.

// UcFirst changes first character to upper case.
// It is based on https://github.com/99designs/gqlgen/blob/master/codegen/templates/templates.go#L205
func UcFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// LcFirst changes first character to lower case.
func LcFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// SnakeToCamel converts snake_case to CamelCase.
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

// RemoveEmpty removes empty string within the slice
func RemoveEmpty(src []string) []string {
	var d []string
	for _, s := range src {
		if s != "" {
			d = append(d, s)
		}
	}
	return d
}
