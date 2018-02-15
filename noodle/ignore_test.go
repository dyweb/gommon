package noodle

import (
	"testing"
	"bytes"

	asst "github.com/stretchr/testify/assert"
	"github.com/dyweb/gommon/util/fsutil"
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
	// TODO: there seems to be no way of checking name and folder pattern ...
	patterns := ignores.Patterns()
	assert.IsType(fsutil.ExactPattern(""), patterns[0])
	assert.IsType(fsutil.WildcardPattern(""), patterns[1])
	assert.IsType(fsutil.WildcardPattern(""), patterns[2])
	assert.IsType(fsutil.WildcardPattern(""), patterns[3])
}

func TestCleanLine(t *testing.T) {
	assert := asst.New(t)
	assert.Equal("", CleanLine("# I should be empty"))
	assert.Equal("", CleanLine(" # I should also be empty"))
	assert.Equal("You can see me", CleanLine("You can see me  # But not the rest"))
	assert.Equal("test/*", CleanLine("test/*"))
}
