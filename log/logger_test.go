package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestLogger_AddField(t *testing.T) {
	assert := asst.New(t)

	logger := NewTestLogger(defaultLevel)
	th := NewTestHandler()
	logger.SetHandler(th)

	logger.AddField(Str("foo", "bar"))
	logger.InfoF("original")
	e, ok := th.getLogByMessage("original")
	assert.True(ok)
	assert.Equal(e.context[0].Key, "foo")
	assert.Equal(e.context[0].Str, "bar")

	// TODO: now AddField is not handling duplication ...
	t.Run("duplication is NOT handled", func(t *testing.T) {
		logger.AddField(Str("foo", "bar2"))
		logger.InfoF("bar2")
		e, ok := th.getLogByMessage("bar2")
		assert.True(ok)
		assert.Equal(e.context[0].Key, "foo")
		assert.Equal(e.context[0].Str, "bar")
		assert.Equal(e.context[1].Key, "foo")
		assert.Equal(e.context[1].Str, "bar2")
	})
}

func TestLogger_WithField(t *testing.T) {
	assert := asst.New(t)

	l1 := NewTestLogger(defaultLevel)
	th := NewTestHandler()
	l1.SetHandler(th)

	l2 := l1.WithField(Str("logger", "l2"))
	l1.AddField(Str("logger", "l1"))

	l1.InfoF("hi1")
	l2.InfoF("hi2")

	e, ok := th.getLogByMessage("hi1")
	assert.True(ok)
	assert.Equal(e.context[0].Str, "l1")

	e, ok = th.getLogByMessage("hi2")
	assert.True(ok)
	assert.Equal(e.context[0].Str, "l2")
}

func TestLogger_WithFields(t *testing.T) {
	assert := asst.New(t)

	logger := NewTestLogger(defaultLevel)
	th := NewTestHandler()
	logger.SetHandler(th)

	logger = logger.WithFields(Int("a", 1), Str("foo", "bar"))
	logger.InfoF("hi")
	e, ok := th.getLogByMessage("hi")
	assert.True(ok)
	assert.Equal(2, len(e.context))
}

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
