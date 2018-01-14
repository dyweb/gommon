package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

// TODO: might use setup and tear down (newer version of go has this built in?)
func TestLogger_SetHandler(t *testing.T) {
	assert := asst.New(t)

	logger := NewFunctionLogger(nil)
	th := NewTestHandler()
	logger.SetHandler(th)
	logger.Info("hello logger")
	assert.True(th.HasLog(InfoLevel, "hello logger"))
	// TODO: test time ...
}

func TestLogger_SetLevel(t *testing.T) {
	assert := asst.New(t)

	logger := NewFunctionLogger(nil)
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
