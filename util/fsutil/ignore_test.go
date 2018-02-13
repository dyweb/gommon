package fsutil

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestWildcardPattern_ShouldIgnore(t *testing.T) {
	t.Skip("wildcard match not implemented")

	assert := asst.New(t)

	t.Run("trailing_star", func(t *testing.T) {
		p := WildcardPattern("abc.*")
		assert.False(p.ShouldIgnore("ab")) // short path, pattern longer than target
		t.Log(p.ShouldIgnore("abc.txt"))
		t.Log(p.ShouldIgnore("abc.t"))
		t.Log(p.ShouldIgnore("abc."))
	})

}
