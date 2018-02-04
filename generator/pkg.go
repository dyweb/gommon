package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dyweb/gommon/config"
	"github.com/dyweb/gommon/util/logutil"
	"github.com/pkg/errors"
)

const generatorName = "gommon"
const generatedFile = "gommon_generated.go"

var log = logutil.NewPackageLogger()

func Generate(root string) error {
	files := Walk(root, DefaultIgnores())
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
	// TODO: config may remove LoadYAMLAsStruct
	if err = config.LoadYAMLAsStruct(file, &cfg); err != nil {
		return errors.WithMessage(err, "can't read config file")
	}
	if rendered, err = cfg.Render(); err != nil {
		return errors.WithMessage(err, "can't render based on config")
	}
	//log.Debugf("%s rendered length %d", file, len(rendered))
	if err = WriteFile(join(dir, generatedFile), rendered); err != nil {
		return errors.WithMessage(err, "can't write rendered file")
	}
	// gotmpl
	if err = cfg.RenderGoTemplate(dir); err != nil {
		return errors.WithMessage(err, "can't render go templates")
	}
	return nil
}
