package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
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

func TestParseLevel(t *testing.T) {
	assert := asst.New(t)
	// strict
	for _, l := range AllLevels {
		s := l.String()
		l2, err := ParseLevel(s, true)
		assert.Nil(err)
		assert.Equal(l, l2)
	}
	// not strict
	levels := []string{"FA", "Pa", "er", "Warn", "infooo", "debugggg", "Tracer"}
	for i := 0; i < len(AllLevels); i++ {
		l, err := ParseLevel(levels[i], false)
		assert.Nil(err)
		assert.Equal(AllLevels[i], l)
	}
	// invalid
	_, err := ParseLevel("haha", false)
	assert.NotNil(err)
}
