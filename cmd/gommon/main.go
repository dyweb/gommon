package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"flag"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/config"
	"github.com/dyweb/gommon/generator"
)

var log = dlog.NewApplicationLogger()
// create new flag set to avoid using default flag
var flags = flag.NewFlagSet("gommon", flag.ExitOnError)
var showHelp = flags.Bool("h", false, "display help")
var verbose = flags.Bool("v", false, "verbose output")
var commands = `
Available Commands:
generate  generate logger methods for struct based on gommon.yml
`

func generate() {
	root := "."
	files := generator.Walk(root, generator.DefaultIgnores)
	hasErr := false
	for _, file := range files {
		dir := filepath.Dir(file)
		segments := strings.Split(dir, string(os.PathSeparator))
		pkg := segments[len(segments)-1]
		c := generator.NewConfig(pkg, file)
		if err := config.LoadYAMLAsStruct(file, &c); err != nil {
			hasErr = true
			log.Warn(err)
			continue
		}
		dst := filepath.Join(dir, "gommon_generated.go")
		f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Warn(err)
			continue
		}
		if rendered, err := c.Render(); err != nil {
			log.Warn(err)
			continue
		} else {
			if _, err := f.Write(rendered); err != nil {
				log.Warn(err)
			}
		}
		f.Close()
		log.Infof("generated %s", dst)
	}
	if hasErr {
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
