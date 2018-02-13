package fsutil

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	assert := asst.New(t)
	//t.Log(Cwd())
	assert.True(Exists("testdata/a.txt"))
	assert.True(Exists("testdata/sub"))
	assert.True(Exists("testdata/sub/a.txt"))
}

func TestFileExists(t *testing.T) {
	assert := asst.New(t)
	assert.True(FileExists("testdata/a.txt"))
	assert.False(FileExists("testdata/sub"))
}

func TestDirExists(t *testing.T) {
	assert := asst.New(t)
	assert.False(DirExists("testdata/a.txt"))
	assert.True(DirExists("testdata/sub"))
}