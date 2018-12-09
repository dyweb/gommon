package log

import (
	"testing"

	"github.com/dyweb/gommon/util/testutil"
	"github.com/stretchr/testify/assert"
)

// test if skip caller works

var lgCallerSkip = NewPackageLogger()

func utilReportError(m string) {
	l := lgCallerSkip.Copy().SetCallerSkip(1)
	l.InfoF(m)
}

func TestLogger_SetCallerSkip(t *testing.T) {
	th := NewTestHandler()
	lgCallerSkip.SetHandler(th)
	lgCallerSkip.EnableSource()

	utilReportError("mie")

	th.HasLog(InfoLevel, "mie")
	l, ok := th.getLogByMessage("mie")
	assert.True(t, ok)
	assert.Equal(t, testutil.GOPATH()+"/src/github.com/dyweb/gommon/log/caller_test.go", l.source.File)
	assert.Equal(t, 24, l.source.Line)
}
