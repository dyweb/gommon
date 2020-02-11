package fsutil

import (
	"bytes"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

var exampleIgnore = `
# example of .noodleignore
# support comment, also blank line should be ignored

vendor # ignore any file or directory whose name is exactly vendor

*.pdf # ignore any file or directory that matches *.pdf

# ignore all the files and directory under test,
# since it also applies to walk, test/sub/example.txt will be ignored as well
# however it is not ignored because match test/* pattern, * does not match separator
# TODO: this is just my assumption ... not tested ...
test/*

# ignore assets/a.partial.html etc.
assets/*.partial.html
`

func TestReadIgnore(t *testing.T) {
	assert := asst.New(t)
	ignores, err := ReadIgnore(bytes.NewReader([]byte(exampleIgnore)))
	assert.Nil(err)
	assert.Equal(4, ignores.Len())
	// TODO: there seems to be no way of checking FileName and folder pattern ...
	patterns := ignores.Patterns()
	assert.IsType(ExactPattern(""), patterns[0])
	assert.IsType(WildcardPattern(""), patterns[1])
	assert.IsType(WildcardPattern(""), patterns[2])
	assert.IsType(WildcardPattern(""), patterns[3])
}

func TestCleanLine(t *testing.T) {
	assert := asst.New(t)
	assert.Equal("", CleanLine("# I should be empty"))
	assert.Equal("", CleanLine(" # I should also be empty"))
	assert.Equal("You can see me", CleanLine("You can see me  # But not the rest"))
	assert.Equal("test/*", CleanLine("test/*"))
}

func TestWildcardPattern_ShouldIgnore(t *testing.T) {
	assert := asst.New(t)

	t.Run("only_star", func(t *testing.T) {
		p := WildcardPattern("*")
		assert.True(p.ShouldIgnore("match anything"))
		assert.False(p.ShouldIgnore("match anything/but slash"))
	})

	t.Run("only_trailing_star", func(t *testing.T) {
		p := WildcardPattern("abc.*")
		assert.False(p.ShouldIgnore("ab")) // short path, pattern longer than target
		assert.True(p.ShouldIgnore("abc.txt"))
		assert.True(p.ShouldIgnore("abc.t"))
		assert.False(p.ShouldIgnore("abc."))
	})

	t.Run("only_question_mark", func(t *testing.T) {
		p := WildcardPattern("?jax_?.html")
		assert.True(p.ShouldIgnore("ajax_a.html"))
	})

	t.Run("only_one_star", func(t *testing.T) {
		p := WildcardPattern("abc*.html")
		assert.True(p.ShouldIgnore("abcd.html"))
		assert.True(p.ShouldIgnore("abcde.html"))
		assert.True(p.ShouldIgnore("abc.ht.html"))
		// should not match /
		assert.False(p.ShouldIgnore("abc/a.html"))
		assert.False(p.ShouldIgnore("abcd/a.html"))
	})

	t.Run("files_under_folder", func(t *testing.T) {
		p := WildcardPattern("assets/*.html")
		assert.True(p.ShouldIgnore("assets/a.html"))
		assert.False(p.ShouldIgnore("assets/a.htm"))
		assert.False(p.ShouldIgnore("assets/a.html/a.html"))
	})

	// TODO: test combination ...
}
