package genutil

import "golang.org/x/tools/imports"

// go.go defines helper when generating go code

// Format calls imports.Process so it behaves like goimports i.e. sort and merge imports.
// Deprecated: call generator.FormatGo
func Format(src []byte) ([]byte, error) {
	// TODO: might need to disable finding import and chose format only
	return imports.Process("", src, nil)
}
