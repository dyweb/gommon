package main

import (
	"fmt"
	"os"
	"path/filepath"

	dlog "github.com/dyweb/gommon/log"
	"github.com/dyweb/gommon/config"
	"github.com/dyweb/gommon/generator"
	"github.com/dyweb/gommon/util"
	"strings"
)

var log = dlog.NewApplicationLogger()

const help = `Gommon utils

Usage:
  Gommon [command]

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
		c := generator.NewConfig(pkg)
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
		f.WriteString(util.GeneratedHeader("gommon", file))
		c.Render(f)
		f.Close()
		log.Infof("generated %s", dst)
	}
	if hasErr {
		os.Exit(1)
	}
}

// TODO: allow testing common gommon features like config, requests, runner
func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, help)
		return
	}
	if os.Args[1] == "generate" {
		generate()
	}
}
