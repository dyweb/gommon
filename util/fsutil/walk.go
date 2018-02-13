package fsutil

import (
	"path/filepath"
	//"os"
)

// TODO: we may want to have the dir know its subdir, which works for Readdir

// I am walking in the sun in around and around (there could be a cycle)
func Walk(root string, ignores Ignores, walkFunc filepath.WalkFunc) []string {
	//filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
	//
	//})
	// TODO: filepath.Walk use Lstat instead of stat, which does not follow symbolic link
	return nil
}
