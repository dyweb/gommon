package linter

import (
	"golang.org/x/tools/imports"
	"testing"
)

func TestImport(t *testing.T)  {
	// for jump into x/tools/imports and x/tools/internal/imports
	imports.Process("foo", nil, nil)
}
