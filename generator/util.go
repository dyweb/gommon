package generator

import (
	"path/filepath"

	"github.com/dyweb/gommon/util/fsutil"
)

func DefaultIgnores() *fsutil.Ignores {
	return fsutil.NewIgnores(
		[]fsutil.IgnorePattern{
			fsutil.ExactPattern(".git"),
			fsutil.ExactPattern("testdata"),
			fsutil.ExactPattern("vendor"),
			fsutil.ExactPattern(".idea"),
			fsutil.ExactPattern(".vscode"),
			fsutil.ExactPattern("legacy"),
		},
		nil,
	)
}

// join is alias for filepath.Join
func join(s ...string) string {
	return filepath.Join(s...)
}

// WriteFile writes file using permission 0664
// NOTE: 0664 is octal literal in Go, the code would compile for 664, but the result file mode is incorrect
// learned this the hard way https://github.com/dyweb/gommon/issues/41
// stat -c %a pkg.go
// NOTE: this code is now in
//func WriteFile(f string, b []byte) error {
//	if err := ioutil.WriteFile(f, b, 0664); err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
