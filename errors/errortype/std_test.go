package errortype_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/errors/errortype"
	"github.com/stretchr/testify/assert"
)

func TestGetRuntimeError(t *testing.T) {
	t.Run("invalid operation", func(t *testing.T) {
		var re error
		var b int
		fn := func() {
			defer func() {
				r := recover()
				re = errors.Wrap(r.(error), "should be divide by zero")
			}()
			// can also be slice index out of range etc.
			a := 1 / b
			t.Logf("a is %d", a)
		}
		b = 0
		fn()
		assert.Equal(t, true, errortype.IsRuntimeError(re))
	})
	t.Run("user panic using string is not runtime.Error", func(t *testing.T) {
		var re error
		fn := func() {
			defer func() {
				r := recover()
				//re = errors.Wrap(r.(error), "user panic")
				// NOTE: it is just a string, it depends on what user pass in
				re = errors.New(r.(string))
			}()
			panic("I just want to say")
		}
		fn()
		assert.Equal(t, false, errortype.IsRuntimeError(re))
	})
	t.Run("user panic using error is just user defined error", func(t *testing.T) {
		var re error
		fn := func() {
			defer func() {
				r := recover()
				re = errors.Wrap(r.(error), "wrap user defined")
			}()
			panic(errors.New("this is my error"))
		}
		fn()
		assert.Equal(t, false, errortype.IsRuntimeError(re))
		assert.Equal(t, "this is my error", errors.Cause(re).Error())
	})
}

func TestGetNetError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, errortype.IsNetError(nil))
	})

	t.Run("invalid url", func(t *testing.T) {
		c := http.Client{Timeout: 200 * time.Millisecond}
		res, err := c.Get("http://a5.me/foo/bar")
		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.True(t, errortype.IsNetError(err))
		// Get http://a5.me/foo/bar: net/http: request canceled (Client.Timeout exceeded while awaiting headers)
		//t.Log(err)
	})

	// TODO: other interesting errors, connection reset by peers etc.
}

func TestGetFsError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, errortype.IsFsError(nil))
	})
	t.Run("open", func(t *testing.T) {
		_, err := os.Open("foo")
		assert.True(t, errortype.IsFsError(err))
		errw := errors.Wrap(err, "wrapped")
		assert.True(t, errortype.IsFsError(errw))
		fsErr, ok := errortype.GetFsError(err)
		assert.True(t, ok)
		assert.Equal(t, err, fsErr)
	})
}
