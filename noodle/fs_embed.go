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
		// FIXME: ignore path does not work because we didn't trim the common prefix
		log.Info(ignores.Patterns())
	}
	fsutil.Walk(root, ignores, func(path string, info os.FileInfo) {
		log.Info(path + "/" + info.Name())
	})
	return nil
}
