package errors_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/testutil"
)

func TestNew(t *testing.T) {
	err := errors.New("don't let me go")
	assert.NotNil(t, err)
	assert.Equal(t, "don't let me go", err.Error())
	assert.Equal(t, "don't let me go", fmt.Sprintf("%v", err))
	terr, ok := err.(errors.Tracer)
	assert.True(t, ok)
	errors.PrintFrames(terr.Stack().Frames())
	assert.Equal(t, 3, len(terr.Stack().Frames()))
}

func TestWrap(t *testing.T) {
	n := errors.Wrap(nil, "nothing")
	assert.Nil(t, n)

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal(t, "can't open closed file: file already closed", fmt.Sprintf("%v", errw))
	terr, ok := errw.(errors.Tracer)
	assert.True(t, ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(t, 3, len(terr.Stack().Frames()))

	errw = errors.Wrap(freshErr(), "wrap again")
	terr, ok = errw.(errors.Tracer)
	assert.True(t, ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(t, 4, len(terr.Stack().Frames()))

	errw = errors.Wrap(wrappedStdErr(), "wrap again")
	terr, ok = errw.(errors.Tracer)
	assert.True(t, ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(t, 4, len(terr.Stack().Frames()))
}

func TestWrapf(t *testing.T) {
	// TODO: need to ensure Wrap and Wrapf are same ...

	// wrap nil return nil
	n := errors.Wrapf(nil, "nothing %d", 2)
	assert.Nil(t, n)

	// wrap standard error attach stack
	errw := errors.Wrapf(os.ErrClosed, "can't open closed file %s", "gommon.yml")
	assert.Equal(t, "can't open closed file gommon.yml: file already closed", fmt.Sprintf("%v", errw))
	terr, ok := errw.(errors.Tracer)
	assert.True(t, ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(t, 3, len(terr.Stack().Frames()))

	// wrap fresh error reuse stack
	ferr := freshErr()
	errw = errors.Wrapf(ferr, "wrap again %d", 2)
	// NOTE: since we changed Stack to interface, we can no longer compare the underlying pointer to struct directly...
	assert.Equal(t, len(ferr.(errors.Tracer).Stack().Frames()), len(errw.(errors.Tracer).Stack().Frames()))

	// wrap wrapped error reuse stack
	errww := errors.Wrapf(errw, "wrap again %d", 3)
	assert.Equal(t, len(errw.(errors.Tracer).Stack().Frames()), len(errww.(errors.Tracer).Stack().Frames()))
}

func TestDirectCause(t *testing.T) {
	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	errww := errors.Wrap(errw, "wrap again")
	assert.Equal(t, os.ErrClosed, errors.Cause(errww))
	assert.Equal(t, os.ErrClosed, errors.Cause(errw))
	assert.NotEqual(t, os.ErrClosed, errors.DirectCause(errww))
	assert.Equal(t, "can't open closed file", errors.DirectCause(errww).(errors.Messenger).Message())
}

func TestIs(t *testing.T) {
	t.Run("flat", func(t *testing.T) {
		assert.False(t, errors.Is(nil, os.ErrClosed), "nil does not match sentinel error")
		assert.True(t, errors.Is(os.ErrClosed, os.ErrClosed), "sentinel errors match themselves")
		assert.False(t, errors.Is(os.ErrNotExist, os.ErrClosed), "different error value should not match")
	})

	t.Run("unwrap wrapper", func(t *testing.T) {
		errw := errors.Wrap(os.ErrClosed, "can't read config")
		assert.True(t, errors.Is(errw, os.ErrClosed), "unwrap wrapped error")

		errww := errors.Wrap(errw, "just wrap again")
		assert.True(t, errors.Is(errww, os.ErrClosed), "unwrap wrapped error")
	})

	t.Run("unwrap multi error", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(os.ErrNotExist)
		merr.Append(os.ErrClosed)
		assert.True(t, errors.Is(merr, os.ErrClosed))
		assert.True(t, errors.Is(merr, os.ErrNotExist))
		assert.False(t, errors.Is(merr, os.ErrNoDeadline))
	})

	t.Run("unwrap wrapper inside multi", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(errors.Wrap(os.ErrNotExist, "can't read config"))
		merr.Append(os.ErrClosed)
		assert.True(t, errors.Is(merr, os.ErrClosed))
		assert.True(t, errors.Is(merr, os.ErrNotExist))
	})

	t.Run("unwrap multi inside wrapper", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(errors.Wrap(os.ErrNotExist, "can't read config"))
		merr.Append(os.ErrClosed)
		errw := errors.Wrap(merr, "just wrap it")
		assert.True(t, errors.Is(errw, os.ErrClosed))
		assert.True(t, errors.Is(errw, os.ErrNotExist))
	})
}

func TestIsType(t *testing.T) {
	t.Run("flat", func(t *testing.T) {
		assert.True(t, errors.IsType(&os.PathError{}, &os.PathError{}))
		assert.False(t, errors.IsType(errors.New("foo"), &os.PathError{}))
	})

	t.Run("unwrap wrapper", func(t *testing.T) {
		errw := errors.Wrap(&os.PathError{Op: "open"}, "can't read config")
		assert.True(t, errors.IsType(errw, &os.PathError{}))
		assert.False(t, errors.IsType(errw, errors.New("foo")))
	})

	t.Run("unwrap error inside multi", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(&os.PathError{Op: "open"})
		merr.Append(os.ErrNotExist)
		assert.True(t, errors.IsType(merr, &os.PathError{}))
	})

	t.Run("unwrap wrapper inside multi", func(t *testing.T) {
		errw := errors.Wrap(&os.PathError{Op: "open"}, "can't read config")
		merr := errors.NewMultiErr()
		merr.Append(errw)
		assert.True(t, errors.IsType(errw, &os.PathError{}))
	})
}

func TestGetType(t *testing.T) {
	t.Run("flat", func(t *testing.T) {
		e := &os.PathError{Op: "open"}
		err, ok := errors.GetType(e, &os.PathError{})
		assert.True(t, ok)
		assert.Equal(t, "open", err.(*os.PathError).Op)
	})

	t.Run("unwrap wrapper", func(t *testing.T) {
		errw := errors.Wrap(&os.PathError{Op: "open"}, "can't read config")
		err, ok := errors.GetType(errw, &os.PathError{})
		assert.True(t, ok)
		assert.Equal(t, "open", err.(*os.PathError).Op)
	})

	t.Run("unwrap error inside multi", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(&os.PathError{Op: "open"})
		err, ok := errors.GetType(merr, &os.PathError{})
		assert.True(t, ok)
		assert.Equal(t, "open", err.(*os.PathError).Op)
	})
}

func TestWrappedError_Message(t *testing.T) {
	msg := "mewo"
	errw := errors.Wrap(os.ErrClosed, msg)
	assert.Equal(t, msg, errw.(errors.Messenger).Message())
}

func ExampleWrap() {
	err := errors.Wrap(os.ErrNotExist, "oops")
	fmt.Println(err)
	// Output:
	// oops: file does not exist
}

// stubs
func freshErr() error {
	return errors.New("I am a fresh error")
}

func wrappedStdErr() error {
	return errors.Wrap(os.ErrClosed, "can't open closed file")
}
