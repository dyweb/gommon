package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestLogger_SetHandler(t *testing.T) {
	assert := asst.New(t)

	logger := NewTestLogger(defaultLevel)
	th := NewTestHandler()
	logger.SetHandler(th)
	logger.Info("hello logger")
	assert.True(th.HasLog(InfoLevel, "hello logger"))
}

func TestLogger_SetLevel(t *testing.T) {
	assert := asst.New(t)

	logger := NewTestLogger(defaultLevel)
	th := NewTestHandler()
	logger.SetHandler(th)
	logger.SetLevel(DebugLevel)
	logger.Info("hi")
	logger.Debug("hi")
	logger.Trace("hi")
	assert.True(th.HasLog(InfoLevel, "hi"))
	assert.True(th.HasLog(DebugLevel, "hi"))
	assert.False(th.HasLog(TraceLevel, "hi"))
}
