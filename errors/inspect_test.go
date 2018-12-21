package errors_test

import (
	"os"
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/stretchr/testify/assert"
)

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
