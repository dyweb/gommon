package linter

import (
	"github.com/dyweb/gommon/linter"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/imports"
	"testing"
)

func TestImport(t *testing.T) {
	// for jump into x/tools/imports and x/tools/internal/imports
	opt := imports.Options{
		TabWidth:   8,
		TabIndent:  true,
		Comments:   true,
		AllErrors:  true,
		FormatOnly: true,
	}
	// This allow extra grouping
	imports.LocalPrefix = "golang.org"
	defer func() {
		imports.LocalPrefix = ""
	}()
	b, err := imports.Process("unordered_import.go", nil, &opt)
	require.Nil(t, err)
	fsutil.WriteFile("unordered_import_goimport.txt", b)
}

func TestGommonImport(t *testing.T) {
	res, err := linter.CheckAndFormatImport("unordered_import.go", linter.GoimportFlags{
		List:       true,
		Diff:       true,
		FormatOnly: true,
	})
	require.Nil(t, err)
	fsutil.WriteFile("unordered_import_gommon.txt", res.Formatted)
}
