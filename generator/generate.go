package generator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/fsutil"
)

// Generate walks all sub directories and generate files based on gommon.yml
func Generate(root string) error {
	var files []string
	// TODO: limit level
	// TODO: allow read ignore from file
	fsutil.Walk(root, DefaultIgnores(), func(path string, info os.FileInfo) {
		//log.Trace(path + "/" + info.Name())
		if info.Name() == GommonConfigFile {
			files = append(files, join(path, info.Name()))
		}
	})
	for _, file := range files {
		if err := GenerateSingle(file); err != nil {
			return err
		}
	}
	return nil
}

// GenerateSingle generates based on a single gommon.yml
func GenerateSingle(file string) error {
	dir := filepath.Dir(file)
	segments := strings.Split(dir, string(os.PathSeparator))
	pkg := segments[len(segments)-1]
	cfg := NewConfigFile(pkg, file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "can't read config file")
	}
	// NOTE: not using Unmarshal strict so new binary still works with old config with deprecated fields
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return errors.Wrap(err, "can't decode config file as YAML")
	}
	if cfg.GoPackage != "" {
		pkg = cfg.GoPackage
	}

	// gommon
	var body bytes.Buffer
	for _, l := range cfg.Loggers {
		b, err := l.Render()
		if err != nil {
			return err
		}
		body.Write(b)
	}
	if body.Len() != 0 {
		var header bytes.Buffer
		header.WriteString(DefaultHeader(file))
		header.WriteString("package " + pkg + "\n\n")
		formatted, err := format.Source(header.Bytes())
		if err != nil {
			return errors.Wrap(err, "error format generated go code")
		}
		if fsutil.WriteFile(join(dir, DefaultGeneratedFile), formatted); err != nil {
			return errors.Wrap(err, "can write generated file to disk")
		}
		log.Debugf("generated %s from %s", join(dir, DefaultGeneratedFile), file)
	}

	// TODO: noodle

	// gotmpl
	for _, tpl := range cfg.GoTemplates {
		if err := tpl.Render(dir); err != nil {
			return err
		}
	}

	// shell
	for _, s := range cfg.Shells {
		if err := s.Render(dir); err != nil {
			return err
		}
	}

	return nil
}
