package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiHandler(t *testing.T) {
	f1, err := os.Create("testdata/f1.log")
	require.Nil(t, err)
	f2, err := os.Create("testdata/f2.log")
	require.Nil(t, err)

	mh := MultiHandler(NewTextHandler(f1), NewTextHandler(f2))
	logger := NewTestLogger(InfoLevel)
	logger.AddFields(Str("s1", "v1"), Int("i1", 1))
	logger.SetHandler(mh)

	logger.Info("should write to both files")
	logger.Warn("this is a warning")
	assert.NoError(t, f1.Close())
	assert.NoError(t, f2.Close())

	b1, err := ioutil.ReadFile("testdata/f1.log")
	require.Nil(t, err)
	b2, err := ioutil.ReadFile("testdata/f2.log")
	require.Nil(t, err)
	assert.Equal(t, b1, b2)
}
