package noodle

import (
	"os"

	"github.com/dyweb/gommon/util/fsutil"
)

// TODO
// walk the folder, keep track of folders

func GenerateEmbed(root string) error {
	var (
		err     error
		ignores *fsutil.Ignores
	)
	if fsutil.FileExists(join(root, ignoreFile)) {
		log.Info("found ignore file")
		if ignores, err = ReadIgnoreFile(join(root, ignoreFile)); err != nil {
			return err
		}
		// set common prefix so ignore path would work
		ignores.SetPathPrefix(root)
		log.Info(ignores.Patterns())
	}
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		log.Info(path + "/" + info.Name())
	})
	return nil
}
