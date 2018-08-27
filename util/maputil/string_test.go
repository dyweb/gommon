package maputil_test

import (
	"testing"

	asst "github.com/stretchr/testify/assert"

	"github.com/dyweb/gommon/util/maputil"
)

func TestCopyStringMap(t *testing.T) {
	assert := asst.New(t)

	oldVal := "old"
	newVal := "new"
	src := map[string]string{
		"foo": oldVal,
	}

	shallowCopy := src // assign map to a new map only create a new reference to same underlying data
	dst := maputil.CopyStringMap(src)
	assert.Equal(len(src), len(dst))
	assert.Equal(src["foo"], dst["foo"])

	src["foo"] = newVal
	assert.NotEqual(newVal, dst["foo"])
	assert.Equal(newVal, shallowCopy["foo"])
}
