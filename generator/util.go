package generator

import (
	"io/ioutil"

	"path/filepath"
)

var DefaultIgnores ignores = []string{"testdata", "vendor", ".idea"}

type ignores []string

func (is *ignores) isIgnored(s string) bool {
	for _, i := range *is {
		if i == s {
			return true
		}
	}
	return false
}

// bfs to find all gommon files
func Walk(root string, ignore ignores) []string {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Warn(err)
	}
	var dirs []string
	var gommonFiles []string
	for _, file := range files {
		name := file.Name()
		if file.IsDir() && !ignore.isIgnored(name) {
			dirs = append(dirs, join(root, name))
			continue
		}
		if name == "gommon.yml" {
			gommonFiles = append(gommonFiles, join(root, name))
		}
	}
	for _, dir := range dirs {
		gommonFiles = append(gommonFiles, Walk(dir, ignore)...)
	}
	return gommonFiles
}

func join(s ...string) string {
	return filepath.Join(s...)
}
