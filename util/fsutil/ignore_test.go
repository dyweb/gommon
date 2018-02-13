package fsutil

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

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
