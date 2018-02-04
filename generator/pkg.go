package generator

import (
	"github.com/dyweb/gommon/config"
	"github.com/dyweb/gommon/util/logutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

const generatorName = "gommon"

var log = logutil.NewPackageLogger()

func Generate(root string) error {
	files := Walk(root, DefaultIgnores())
	for _, file := range files {
		if dst, err := GenerateSingle(file); err != nil {
			return err
		} else {
			log.Infof("generated %s", dst)
		}
	}
	return nil
}

func GenerateSingle(file string) (string, error) {
	var (
		f        *os.File
		err      error
		rendered []byte
		dst      string
	)

	dir := filepath.Dir(file)
	segments := strings.Split(dir, string(os.PathSeparator))
	pkg := segments[len(segments)-1]
	cfg := NewConfig(pkg, file)
	// TODO: config may replace LoadYAMLAsStruct
	if err = config.LoadYAMLAsStruct(file, &cfg); err != nil {
		return dst, errors.WithMessage(err, "can't read config file")
	}
	// TODO: this does not apply for gotmpl, which can specify their own src and destination
	dst = filepath.Join(dir, "gommon_generated.go")
	if f, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		return dst, errors.WithMessage(err, "can't create file for write")
	}
	defer f.Close()
	if rendered, err = cfg.Render(); err != nil {
		return dst, errors.WithMessage(err, "can't render based on config")
	} else {
		if _, err = f.Write(rendered); err != nil {
			return dst, errors.Wrap(err, "can't write rendered data to file")
		}
	}
	return dst, nil
}
