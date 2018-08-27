package generator

import (
	"bytes"
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
	var buf bytes.Buffer
	for _, l := range cfg.Loggers {
		b, err := l.RenderBody(dir)
		if err != nil {
			return err
		}
		buf.Write(b)
	}
	if buf.Len() != 0 {
		// TODO: have the imports
		writeGoFile(pkg, nil, buf.Bytes(), file, join(dir, DefaultGeneratedFile))
	}

	// TODO: write package, imports and write file

	// TODO: noodle
	buf.Reset()

	// gotmpl
	for _, tpl := range cfg.GoTemplates {
		if tpl.IsGo() {
			b, err := tpl.RenderBody(dir)
			if err != nil {
				return err
			}
			// TODO: write imports (NO package) and write file
		} else {
			if err := tpl.Render(dir); err != nil {
				return err
			}
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
