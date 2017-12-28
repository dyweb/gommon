package log

import (
	"testing"

	"bytes"
	asst "github.com/stretchr/testify/assert"
	"io"
	"os"
)

func TestLogger_AddFilter(t *testing.T) {
	assert := asst.New(t)
	logger := NewLogger()
	allow := make(map[string]bool)
	pkgFilter := NewPkgFilter(allow)
	logger.AddFilter(pkgFilter, DebugLevel)
	assert.Equal(1, len(logger.Filters[DebugLevel]))
}

func TestLogger_NewEntryWithPkg(t *testing.T) {
	assert := asst.New(t)
	logger := NewLogger()
	entry := logger.NewEntryWithPkg("x.dummy")
	assert.Equal(1, len(entry.Fields))
	entry.Info("show me the pkg")
}

func TestLogger_SetLevel(t *testing.T) {
	assert := asst.New(t)

	var b bytes.Buffer
	writer := io.MultiWriter(os.Stdout, &b)
	logger := NewLogger()
	logger.Out = writer
	entry1 := logger.NewEntryWithPkg("pkg1")
	assert.Equal(entry1.EntryLevel, InfoLevel)
	entry2 := logger.NewEntryWithPkg("pkg2")

	entry1.Info("should see me")
	entry1.Debug("should not see me")
	assert.NotContains(b.String(), "not")

	assert.Nil(entry2.SetEntryLevel("trace"))
	assert.Nil(logger.SetLevel("debug"))

	entry1.Debug("should see me")
	entry1.Trace("should not see me")
	entry2.Trace("should see me")
	assert.NotContains(b.String(), "not")
}
