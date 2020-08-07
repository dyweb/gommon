package linter

import (
	"bytes"
	"fmt"
	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"go/ast"
	"golang.org/x/tools/imports"
	"io/ioutil"
)

// import.go checks if there are deprecated import and sort import by group

// GoimportFlags is a struct whose fields map to goimports command flags.
type GoimportFlags struct {
	List  bool
	Write bool
	Diff  bool
	// TODO: srcdir, single file as if in dir xxx
	AllErrors   bool
	LocalPrefix string
	FormatOnly  bool
}

func CheckAndFormatFile(p string, flags GoimportFlags) error {
	opt := &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
		// NOTE: we don't have Env because it is in internal/imports and relies on default env.
		AllErrors:  flags.AllErrors,
		FormatOnly: flags.FormatOnly,
	}
	// TODO: we can't set LocalPrefix in env, and the package var seems to be left for this purpose
	if flags.LocalPrefix != "" {
		imports.LocalPrefix = flags.LocalPrefix
	}

	log.Debugf("check and format %s", p)
	src, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	goimportRes, err := imports.Process(p, src, opt)
	if err != nil {
		return errors.Wrap(err, "error calling goimports")
	}

	// TODO: parse and check and format

	// NOTE: Copied from processFile in goimports
	res := goimportRes
	if !bytes.Equal(src, res) {
		if flags.List {
			fmt.Println(p)
		}
		if flags.Write {
			// TODO: why goimports use 0 for perm ...
			if err := fsutil.WriteFile(p, res); err != nil {
				return err
			}
		}
		if flags.Diff {

		}
	}
	return nil
}

func CheckImport(f *ast.File) error {
	return nil
}

func FormatImport(f *ast.File) error {
	return nil
}
