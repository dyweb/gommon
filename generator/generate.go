package generator

import (
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
	var (
		err      error
		rendered []byte
	)

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

	// gommon, logger
	if rendered, err = cfg.RenderGommon(); err != nil {
		return errors.Wrap(err, "can't render based on config")
	}
	if len(rendered) != 0 {
		//log.Debugf("%s rendered length %d", file, len(rendered))
		if err = fsutil.WriteFile(join(dir, GeneratedFile), rendered); err != nil {
			return errors.Wrap(err, "can't write rendered gommon file")
		}
		log.Debugf("generated %s from %s", join(dir, GeneratedFile), file)
	} else {
		// FIXME: (at15) this log is not accurate, gommon will have more than just logger identity
		log.Debugf("%s does not have gommon logger config", dir)
	}

	return nil
}
