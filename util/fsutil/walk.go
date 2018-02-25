package fsutil

import (
	"io/ioutil"
	"os"

	"github.com/dyweb/gommon/errors"
)

type WalkFunc func(path string, info os.FileInfo)

// Walk traverse the directory with ignore patterns in Pre-Order DFS
func Walk(root string, ignores *Ignores, walkFunc WalkFunc) error {
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
		walkFunc(root, file)
		if file.IsDir() {
			if err := Walk(path, ignores, walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}
