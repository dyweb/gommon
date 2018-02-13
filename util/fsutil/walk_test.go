package fsutil

import (
	"testing"
	"os"

	asst "github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	assert := asst.New(t)

	ls := func(ignores *Ignores) []string {
		var files []string
		err := Walk("testdata", ignores, func(path string, info os.FileInfo) {
			files = append(files, join(path, info.Name()))
		})
		assert.Nil(err)
		return files
	}

	t.Run("no_ignore", func(t *testing.T) {
		assert.Equal([]string{
			"testdata/a.txt",
			"testdata/b.txt",
			"testdata/sub",
			"testdata/sub/a.txt",
			"testdata/sub/b.txt",
			"testdata/sub2",
			"testdata/sub2/a.txt",
			"testdata/sub2/b.txt",
		}, ls(nil))
	})

	t.Run("ignore_single_file", func(t *testing.T) {
		ignores := NewIgnores(
			[]IgnorePattern{
				ExactPattern("a.txt"),
			},
			nil,
		)
		assert.Equal([]string{
			"testdata/b.txt",
			"testdata/sub",
			"testdata/sub/b.txt",
			"testdata/sub2",
			"testdata/sub2/b.txt",
		}, ls(ignores))
	})

	t.Run("ignore_single_path", func(t *testing.T) {
		ignores := NewIgnores(
			nil,
			[]IgnorePattern{
				// TODO: we have to include the real root in the pattern ... which is quite counter intuitive
				WildcardPattern("testdata/sub2/*.txt"),
			},
		)
		assert.Equal([]string{
			"testdata/a.txt",
			"testdata/b.txt",
			"testdata/sub",
			"testdata/sub/a.txt",
			"testdata/sub/b.txt",
			"testdata/sub2",
		}, ls(ignores))
	})

	t.Run("ignore_name_and_path", func(t *testing.T) {
		ignores := NewIgnores(
			[]IgnorePattern{
				ExactPattern("a.txt"),
			},
			[]IgnorePattern{
				// TODO: we have to include the real root in the pattern ... which is quite counter intuitive
				WildcardPattern("testdata/sub2/*.txt"),
			},
		)
		assert.Equal([]string{
			"testdata/b.txt",
			"testdata/sub",
			"testdata/sub/b.txt",
			"testdata/sub2",
		}, ls(ignores))
	})
}
