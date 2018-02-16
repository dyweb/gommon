package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dyweb/gommon/config"
	"github.com/dyweb/gommon/util/fsutil"
	"github.com/pkg/errors"
)

func Generate(root string) error {
	var files []string
	// TODO: limit level
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
	log.Debugf("generated %s from %s", join(dir, generatedFile), file)
	// gotmpl
	if err = cfg.RenderGoTemplate(dir); err != nil {
		return errors.WithMessage(err, "can't render go templates")
	}
	return nil
}
