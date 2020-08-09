// gommon is the commandline util for generator
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dyweb/gommon/linter"
	"github.com/spf13/cobra"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/generator"
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/cli"
	"github.com/dyweb/gommon/noodle"
	"github.com/dyweb/gommon/util/fsutil"
)

var logReg = dlog.NewRegistry()
var log = logReg.Logger()

var verbose = false
var (
	version   string
	commit    string
	buildTime string
	buildUser string
	goVersion = runtime.Version()
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gommon",
		Short: "gommon helpers",
		Long:  "Generate go files for gommon",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				dlog.SetLevel(dlog.DebugLevel)
				dlog.EnableSource()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	// global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	// ver
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  "Print current version " + version,
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Printf("version: %s\n", version)
				fmt.Printf("commit: %s\n", commit)
				fmt.Printf("build time: %s\n", buildTime)
				fmt.Printf("build user: %s\n", buildUser)
				fmt.Printf("go version: %s\n", goVersion)
			} else {
				fmt.Println(version)
			}
		},
	}
	// sub commands
	rootCmd.AddCommand(
		versionCmd,
		genCmd(),
		addBuildIgnoreCmd(),
		formatCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// triggers generator for logger, go template and embedding asset
func genCmd() *cobra.Command {
	gen := cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "generate code based on gommon.yml",
		Run: func(cmd *cobra.Command, args []string) {
			root := "."
			if err := generator.Generate(root); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	root := ""
	name := ""
	pkg := ""
	output := ""
	// TODO: might put noodle as its own top level command?
	noodleCmd := &cobra.Command{
		Use: "noodle",
		Short: "bundle static assets in one directory as a single go file, " +
			"it does not support bundle multiple file into one directory",
		Example: "gommon generate noodle --root assets --output gen/noodle.go --pkg gen --name Assets",
		Run: func(cmd *cobra.Command, args []string) {
			checks := map[string]string{
				"root":   root,
				"name":   name,
				"pkg":    pkg,
				"output": output,
			}
			merr := errors.NewMultiErr()
			for k, v := range checks {
				if v == "" {
					merr.Append(errors.Errorf("%s is required", k))
				}
			}
			if merr.HasError() {
				log.Fatal(merr.Error())
				return
			}
			cfg := noodle.EmbedConfig{
				Src:     root,
				Name:    name,
				Dst:     output,
				Package: pkg,
			}
			if err := noodle.GenerateEmbedFile(cfg); err != nil {
				log.Fatal(err)
				return
			}
			log.Infof("generated %s from %s with package %s and name %s", output, root, pkg, name)
		},
	}
	noodleCmd.Flags().StringVar(&root, "root", "", "path of assets folder")
	noodleCmd.Flags().StringVar(&name, "name", "Asset", "name of generate ")
	noodleCmd.Flags().StringVar(&pkg, "pkg", "gen", "go package of generated file")
	noodleCmd.Flags().StringVar(&output, "output", "noodle.go", "path for generated file")
	gen.AddCommand(noodleCmd)
	return &gen
}

// add // +build ignore to files before moving them to legacy folder
func addBuildIgnoreCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add-build-ignore",
		Short: "Add // +build ignore to go files",
		Run: func(cmd *cobra.Command, args []string) {
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("error get current directory: %s", err)
				return
			}
			buildIgnore := "// +build ignore\n\n"
			ignores := fsutil.NewIgnores([]fsutil.IgnorePattern{
				fsutil.ExactPattern(".git"),
				fsutil.ExactPattern("testdata"),
				fsutil.ExactPattern("vendor"),
				fsutil.ExactPattern(".idea"),
				fsutil.ExactPattern(".vscode"),
			}, nil)
			for _, p := range args {
				err = fsutil.Walk(filepath.Join(pwd, p), ignores, func(path string, info os.FileInfo) {
					if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
						return
					}
					f := filepath.Join(path, info.Name())
					b, err := ioutil.ReadFile(f)
					if err != nil {
						log.Fatalf("error read file %s", err)
						return
					}
					// NOTE: we only
					prefix := []byte(buildIgnore)
					if bytes.HasPrefix(b, prefix) {
						log.Warnf("%s already have build ignore prefix", f)
						return
					}
					// prepend prefix
					b = append(prefix, b...)
					if err := fsutil.WriteFile(f, b); err != nil {
						log.Fatalf("error write file with build prefix: %s", err)
						return
					}
					log.Infof("updated %s", f)
				})
				if err != nil {
					log.Fatalf("error walk %s", err)
					return
				}
			}
		},
	}
	return &cmd
}

// gommon format
func formatCmd() *cobra.Command {
	var flags linter.GoimportFlags
	processFile := func(path string, info os.FileInfo, err error) error {
		if err == nil && fsutil.IsGoFile(info) {
			return linter.CheckAndFormatImportToStdout(path, flags)
		}
		// Skip directory and stop on walk error
		return err
	}

	run := func(paths []string) error {
		for _, p := range paths {
			switch dir, err := os.Stat(p); {
			case err != nil:
				return err
			case dir.IsDir():
				// TODO: walk w/ ignore like generator
				if err := filepath.Walk(p, processFile); err != nil {
					return err
				}
			default:
				if err := processFile(p, dir, nil); err != nil {
					return err
				}
			}
		}
		return nil
	}

	cmd := cobra.Command{
		Use:   "format",
		Short: "Format go code like goimports with custom rules",
		Run: func(cmd *cobra.Command, args []string) {
			paths := args
			if len(paths) == 0 {
				log.Fatal("format stdin is not implemented")
				return
			}
			if err := run(paths); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().BoolVarP(&flags.List, "list", "l", false, "list files whose formatting differs from goimports")
	cmd.Flags().BoolVarP(&flags.Write, "write", "w", false, "write result to (source) file instead of stdout, i.e. in place update")
	cmd.Flags().BoolVarP(&flags.Diff, "diff", "d", false, "display diffs instead of rewriting files")
	cmd.Flags().BoolVarP(&flags.AllErrors, "errors", "e", false, "report all errors (not just the first 10 on different lines)")
	cmd.Flags().StringVar(&flags.LocalPrefix, "local", "", "put imports beginning with this string after 3rd-party packages; comma-separated list")
	cmd.Flags().BoolVar(&flags.FormatOnly, "format-only", false, "if true, don't fix imports and only format. In this mode, goimports is effectively gofmt, with the addition that imports are grouped into sections.")
	return &cmd
}

func init() {
	dlog.SetHandler(cli.New(os.Stderr, true))
}
