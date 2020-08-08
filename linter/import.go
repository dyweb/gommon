package linter

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
	"golang.org/x/tools/imports"
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

// CheckAndFormatImport calls CheckAndFormatImport and prints to stdout
func CheckAndFormatImport(p string, flags GoimportFlags) error {
	return CheckAndFormatImportTo(os.Stdout, p, flags)
}

func CheckAndFormatImportTo(out io.Writer, p string, flags GoimportFlags) error {
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

	// TODO: call customized check and format

	// NOTE: Copied from processFile in goimports
	res := goimportRes
	// There is diff after format, try update file in place and print diff.
	if !bytes.Equal(src, res) {
		// Print file name
		if flags.List {
			fmt.Fprintln(out, p)
		}
		// Update file directly
		if flags.Write {
			// TODO: why goimports use 0 for file permission?
			if err := fsutil.WriteFile(p, res); err != nil {
				return err
			}
		}
		// Shell out to diff
		if flags.Diff {
			// TODO(upstream): goimports is using Printf instead Fprintf
			fmt.Fprintf(out, "diff -u %s %s\n", filepath.ToSlash(p+".orig"), filepath.ToSlash(p))
			diff, err := diffBytes(src, res, p)
			if err != nil {
				log.Warnf("diff failed %s", err)
			} else {
				out.Write(diff)
			}
		}
	}

	// No flags, dump formatted (may or may not changed) to stdout
	if !flags.List && !flags.Write && !flags.Diff {
		if _, err = out.Write(res); err != nil {
			return err
		}
	}

	return nil
}

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
	return nil
}

func FormatImport(f *ast.File) error {
	return nil
}
