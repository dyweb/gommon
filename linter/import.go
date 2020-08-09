package linter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"golang.org/x/tools/imports"
)

// import.go checks if there are deprecated import and sort import by group

const TabWidth = 8

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

// CheckAndFormatImportToStdout calls CheckAndFormatImport and update file and/or print diff
func CheckAndFormatImportToStdout(p string, flags GoimportFlags) error {
	res, err := CheckAndFormatImport(p, flags)
	if err != nil {
		return err
	}
	return printImportDiff(os.Stdout, flags, res.Path, res.Src, res.Formatted)
}

type FormatResult struct {
	Path      string // file path
	Src       []byte
	Formatted []byte
}

func CheckAndFormatImport(p string, flags GoimportFlags) (*FormatResult, error) {
	log.Debugf("check and format %s", p)
	src, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	// goimports
	opt := imports.Options{
		TabWidth:  TabWidth,
		TabIndent: true,
		Comments:  true,
		Fragment:  false, // Set it to false because we are reading a full go file
		// NOTE: we don't have Env because it is in internal/imports and relies on default env.
		AllErrors:  flags.AllErrors,
		FormatOnly: flags.FormatOnly,
	}
	// TODO: we can't set LocalPrefix in env, and the package var seems to be left for this purpose
	if flags.LocalPrefix != "" {
		imports.LocalPrefix = flags.LocalPrefix
	}
	goimportRes, err := imports.Process(p, src, &opt)
	if err != nil {
		return nil, errors.Wrap(err, "error calling goimports")
	}

	// Parse again ...
	fset := token.NewFileSet()
	// NOTE: we can't use parser.ImportOnly because we need to print ast back to file again.
	mode := parser.ParseComments
	if flags.AllErrors {
		mode |= parser.AllErrors
	}
	f, err := parser.ParseFile(fset, p, goimportRes, mode)
	if err != nil {
		return nil, errors.Wrap(err, "error parse goimport formatted code")
	}

	// Check before format
	// TODO: pass in lint rule
	if err := CheckImport(f); err != nil {
		return nil, err
	}

	// Format again with custom rules
	// TODO: pass in format rule
	res, err := FormatImport(f, fset)
	if err != nil {
		return nil, err
	}

	return &FormatResult{
		Path:      p,
		Src:       src,
		Formatted: res,
	}, nil
}

// NOTE: Copied from processFile in goimports
// If there is diff after format, try update file in place and print diff.
func printImportDiff(out io.Writer, flags GoimportFlags, p string, src []byte, formatted []byte) error {
	if bytes.Equal(src, formatted) {
		return nil
	}

	// Print file name
	if flags.List {
		fmt.Fprintln(out, p)
	}
	// Update file directly
	if flags.Write {
		// TODO: why goimports use 0 for file permission?
		if err := fsutil.WriteFile(p, formatted); err != nil {
			return err
		}
	}
	// Shell out to diff
	if flags.Diff {
		// TODO(upstream): goimports is using Printf instead Fprintf
		fmt.Fprintf(out, "diff -u %s %s\n", filepath.ToSlash(p+".orig"), filepath.ToSlash(p))
		diff, err := diffBytes(src, formatted, p)
		if err != nil {
			log.Warnf("diff failed %s", err)
		} else {
			out.Write(diff)
		}
	}

	// No flags, dump formatted (may or may not changed) to stdout
	if !flags.List && !flags.Write && !flags.Diff {
		if _, err := out.Write(formatted); err != nil {
			return err
		}
	}
	return nil
}

// Writes bytes to two temp files and shell out to diff.
// Replaces file name in output
func diffBytes(a []byte, b []byte, p string) ([]byte, error) {
	files, err := fsutil.WriteTempFiles("", "gomfmt", a, b)
	if err != nil {
		return nil, err
	}
	defer fsutil.RemoveFiles(files)

	diff, err := fsutil.Diff(files[0], files[1])
	if err != nil {
		return nil, err
	}
	// No diff
	if len(diff) == 0 {
		return nil, nil
	}

	// Replace temp file name with original file path
	// NOTE: it is based on replaceTempFilename in goimports
	segs := bytes.SplitN(diff, []byte{'\n'}, 3)
	// NOTE: we ignore invalid diff output and returns whatever the output is.
	// This is different from goimports which stops and return nil.
	if len(segs) < 3 {
		return diff, nil
	}

	// COPY: Copied from goimports replaceTempFilename
	// Preserve timestamps.
	var t0, t1 []byte
	if i := bytes.LastIndexByte(segs[0], '\t'); i != -1 {
		t0 = segs[0][i:]
	}
	if i := bytes.LastIndexByte(segs[1], '\t'); i != -1 {
		t1 = segs[1][i:]
	}
	// Always print filepath with slash separator.
	f := filepath.ToSlash(p)
	segs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	segs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))
	return bytes.Join(segs, []byte{'\n'}), nil
}

func CheckImport(f *ast.File) error {
	// TODO: check the import
	return nil
}

func FormatImport(f *ast.File, fset *token.FileSet) ([]byte, error) {
	// TODO: adjust import group before print
	cfg := printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: TabWidth,
		Indent:   0,
	}
	var buf bytes.Buffer
	if err := cfg.Fprint(&buf, fset, f); err != nil {
		return nil, errors.Wrap(err, "error print ast")
	}
	return buf.Bytes(), nil
}
