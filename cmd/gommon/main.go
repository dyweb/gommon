// gommon is the commandline util for generator
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/generator"
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/cli"
	"github.com/dyweb/gommon/noodle"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/dyweb/gommon/util/logutil"
)

var log = dlog.NewApplicationLogger()
var verbose = false
var (
	version   string
	commit    string
	buildTime string
	buildUser string
	goVersion = runtime.Version()
)

func main() {
	// TODO: most code here are copied from go.ice's cli package, dependency management might break if we import go.ice which also import gommon
	rootCmd := &cobra.Command{
		Use:   "gommon",
		Short: "gommon helpers",
		Long:  "Generate go files for gommon",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				dlog.SetLevelRecursive(log, dlog.DebugLevel)
				dlog.EnableSourceRecursive(log)
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
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genCmd() *cobra.Command {
	gen := &cobra.Command{
		Use:   "generate",
		Short: "generate code based on gommon.yml",
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
		Use:   "noodle",
		Short: "bundle static assets as a single go file",
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
				Root: root,
				Name: name,
			}
			b, err := noodle.GenerateEmbeds([]noodle.EmbedConfig{cfg}, pkg)
			if err != nil {
				log.Fatal(err)
				return
			}
			if err := fsutil.WriteFile(output, b); err != nil {
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
	return gen
}

func init() {
	log.AddChild(logutil.Registry)
	dlog.SetHandlerRecursive(log, cli.New(os.Stderr, true))
}
