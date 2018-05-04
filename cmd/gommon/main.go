// gommon is the commandline util for generator
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/dyweb/gommon/generator"
	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/log/handlers/cli"
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

// TODO: allow testing common gommon features like config, requests, runner
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
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
	genCmd := &cobra.Command{
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
	rootCmd.AddCommand(versionCmd, genCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	log.AddChild(logutil.Registry)
	dlog.SetHandlerRecursive(log, cli.New(os.Stderr, true))
}
