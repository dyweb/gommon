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
	"strings"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"golang.org/x/tools/imports"
)

// import.go checks if there are deprecated import and sort import by group

// ----------------------------------------------------------------------------
// goimports shim

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
	return printImportDiff(os.Stdout, flags, p, res.Src, res.Formatted)
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
	// NOTE: we can't use parser.ImportOnly because we need all the AST to print the file back.
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
	if err := CheckImport(f, nil); err != nil {
		return nil, err
	}

	// Format again with custom rules
	// TODO: pass in format config from outside
	fmtCfg := ImportFormatConfig{}
	res, err := FormatImport(f, fset, fmtCfg)
	if err != nil {
		return nil, err
	}

	return &FormatResult{
		Src:       src,
		Formatted: res,
	}, nil
}

// ----------------------------------------------------------------------------
// Import Context

// ImportContext contains information from ImportSpec and current file.
type ImportContext struct {
	Name    string // import alias or dot import
	Path    string // import path e.g. github.com/foo/bar/bla
	File    string // file name only, e.g. foo.go
	Folder  string // folder, e.g. github.com/gommon/errors TODO: this depends on where the command runs
	Package string // package of current file
}

func importSpecToContext(spec *ast.ImportSpec, pkg string) ImportContext {
	return ImportContext{
		Name: nameIfNotNil(spec.Name),
		Path: spec.Path.Value,
		// TODO: file etc.
		Package: pkg,
	}
}

// ----------------------------------------------------------------------------
// Check Import

// TODO: convert rules from config?
type ImportCheckConfig struct {
}

type ImportCheckRule interface {
	RuleName() string
	Check(ctx ImportContext) error
}

type ImportDeprecated struct {
	Path      string // import path of deprecated package e.g. github.com/foo/bar
	Suggested string // suggested new package e.g. github.com/bar/foo
	Reason    string // why, link to issues etc.
}

func CheckImport(f *ast.File, rules []ImportCheckRule) error {
	pkg := nameIfNotNil(f.Name)
	imps, _ := extractImports(f)
	merr := errors.NewMulti()
	for _, imp := range imps {
		for _, r := range rules {
			// TODO: need to attach more context and maybe have extra grouping
			merr.Append(r.Check(importSpecToContext(imp, pkg)))
		}
	}
	return merr.ErrorOrNil()
}

// ----------------------------------------------------------------------------
// Format Import

// FormatResult avoids remembering order of src and dst when returning two bytes.
type FormatResult struct {
	Src       []byte
	Formatted []byte
}

// TODO: missing other rules
type ImportFormatConfig struct {
	// Path of current project, imports from current project are grouped together
	ProjectPath string

	GroupRules []ImportGroupRule
}

// ImportGroupRule checks if an import belongs to a group.
// TODO: might change to interface
type ImportGroupRule struct {
	// Name is a short string for enabling/disabling rules.
	Name string
	// Match determines if the import spec matches this group.
	Match func(ctx ImportContext) bool
}

var importStd = ImportGroupRule{
	Name: "std",
	Match: func(ctx ImportContext) bool {
		// NOTE: it is based on goimports, if the path has no domain name, it is standard library.
		if !strings.Contains(ctx.Path, ".") {
			return true
		}
		return false
	},
}

var importCurrentPackage = ImportGroupRule{
	Name: "currentPackage",
	Match: func(ctx ImportContext) bool {
		// TODO: this does not work if we don't run gommon format at project root
		if strings.HasSuffix(ctx.Folder, ctx.Path) {
			return true
		}
		// FIXME: the package name may not match the folder name foo/bar, package boar
		return false
	},
}

func DefaultImportGroupRules() []ImportGroupRule {
	return []ImportGroupRule{
		importCurrentPackage,
		importStd,
	}
}

// TODO: split the API, one use cfg, one using rules (i.e. full customization)
func FormatImport(f *ast.File, fset *token.FileSet, cfg ImportFormatConfig) ([]byte, error) {
	rules := cfg.GroupRules
	if len(rules) == 0 {
		// TODO: need to pass project dir, and make it easy for user to define rules
		rules = DefaultImportGroupRules()
	}
	// Change import order and do the grouping.
	if err := adjustImport(f, fset, rules); err != nil {
		return nil, err
	}

	// TODO: apply the rules
	// TODO: the hard part is how to keep the position valid after modifying AST

	// Print AST
	pCfg := printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: TabWidth,
		Indent:   0,
	}
	var buf bytes.Buffer
	if err := pCfg.Fprint(&buf, fset, f); err != nil {
		return nil, errors.Wrap(err, "error print ast")
	}
	return buf.Bytes(), nil
}

func adjustImport(f *ast.File, fset *token.FileSet, rules []ImportGroupRule) error {
	specs, nBlocks := extractImports(f)
	if nBlocks > 1 {
		return errors.Errorf("goimports not called, expect import declaration merged into one block, but got %d", nBlocks)
	}
	// TODO: linter says decl.Specs can be nil ... writing a linter w/ linter error
	for _, imp := range specs {
		log.Debugf("adjustImport: spec %s %s", nameIfNotNil(imp.Name), imp.Path.Value)
	}
	return nil
}

func extractImports(f *ast.File) ([]*ast.ImportSpec, int) {
	var (
		imps    []*ast.ImportSpec
		nBlocks int
	)
	for _, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			break
		}
		nBlocks++
		for _, spec := range d.Specs {
			imps = append(imps, spec.(*ast.ImportSpec))
		}
	}
	return imps, nBlocks
}

func nameIfNotNil(id *ast.Ident) string {
	if id == nil {
		return ""
	}
	return id.Name
}

// ----------------------------------------------------------------------------
// Diff formatted go file

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
// TODO: we can use diffBytes in other packages, might move it to a util package, bytesutil?
func diffBytes(a []byte, b []byte, filename string) ([]byte, error) {
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
	f := filepath.ToSlash(filename)
	segs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	segs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))
	return bytes.Join(segs, []byte{'\n'}), nil
}
