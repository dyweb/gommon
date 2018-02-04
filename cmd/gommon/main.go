package main

import (
	"fmt"
	"os"
	"flag"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/generator"
	"github.com/dyweb/gommon/util/logutil"
)

var log = dlog.NewApplicationLogger()
// create new flag set to avoid using default flag
var flags = flag.NewFlagSet("gommon", flag.ExitOnError)
var showHelp = flags.Bool("h", false, "display help")
var verbose = flags.Bool("v", false, "verbose output")
var commands = `fi
Available Commands:

generate  generate logger methods for struct based on gommon.yml
`

func generate() {
	root := "."
	if err := generator.Generate(root); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func help() {
	flags.Usage()
	fmt.Fprint(os.Stderr, commands)
	os.Exit(1)
}

func parseFlags(args []string) {
	if err := flags.Parse(args); err != nil {
		log.Error(err)
	}
	if *showHelp {
		help()
	}
	if *verbose {
		dlog.SetLevelRecursive(log, dlog.DebugLevel)
	}
}

// TODO: allow testing common gommon features like config, requests, runner
func main() {
	if len(os.Args) < 2 {
		help()
	}
	parseFlags(os.Args[1:])
	if os.Args[1] == "generate" {
		parseFlags(os.Args[2:])
		generate()
	}
}

func init() {
	log.AddChild(logutil.Registry)
}
