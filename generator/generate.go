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

func Generate(root string) error {
	var files []string
	// TODO: limit level
	// TODO: allow read ignore from file
	fsutil.Walk(root, DefaultIgnores(), func(path string, info os.FileInfo) {
		//log.Trace(path + "/" + info.Name())
		if info.Name() == configFile {
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

func GenerateSingle(file string) error {
	var (
		err      error
		rendered []byte
	)

	dir := filepath.Dir(file)
	segments := strings.Split(dir, string(os.PathSeparator))
	pkg := segments[len(segments)-1]
	cfg := NewConfig(pkg, file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "can't read config file")
	}
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return errors.Wrap(err, "can't decode config file as YAML")
	}

	// gommon logger
	if rendered, err = cfg.RenderGommon(); err != nil {
		return errors.Wrap(err, "can't render based on config")
	}
	if len(rendered) != 0 {
		//log.Debugf("%s rendered length %d", file, len(rendered))
		if err = fsutil.WriteFile(join(dir, generatedFile), rendered); err != nil {
			return errors.Wrap(err, "can't write rendered gommon file")
		}
		log.Debugf("generated %s from %s", join(dir, generatedFile), file)
	} else {
		// FIXME: (at15) this log is not accurate, gommon will have more than just logger identity
		log.Debugf("%s does not have gommon logger config", dir)
	}

	// gotmpl
	if err = cfg.RenderGoTemplate(dir); err != nil {
		return errors.Wrap(err, "can't render go templates")
	}

	// logger
	if err = cfg.RenderShell(dir); err != nil {
		return errors.Wrap(err, "can't render using shell commands")
	}

	return nil
}
