package fsutil

import (
	"io/ioutil"
	"os"

	"github.com/dyweb/gommon/errors"
)

// TODO: should support skip dir like filepath.Walk
type WalkFunc func(path string, info os.FileInfo)

// Walk traverse the directory with ignore patterns in Pre-Order DFS
func Walk(root string, ignores *Ignores, walkFunc WalkFunc) error {
	// TODO: validate ignores or assign a default accept all
	if ignores == nil {
		ignores = AcceptAll
	}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return errors.Wrapf(err, "can't read dir %s", root)
	}
	for _, file := range files {
		if ignores.IgnoreName(file.Name()) {
			continue
		}
		path := join(root, file.Name())
		if ignores.IgnorePath(path) {
			continue
		}
		// TODO: this path behavior is also different filepath.Walk
		walkFunc(root, file)
		if file.IsDir() {
			if err := Walk(path, ignores, walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}
