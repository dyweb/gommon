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

func TestWalk(t *testing.T) {
	t.Run("stop", func(t *testing.T) {
		merr := errors.NewMultiErr()
		merr.Append(os.ErrNotExist)
		merr.Append(os.ErrNotExist)

		var errs []error
		errors.Walk(merr, func(err error) (stop bool) {
			errs = append(errs, err)
			return true
		})
		assert.Equal(t, 1, len(errs), "WalkFunc is only called once because it returns true for stop")

		// https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
		errs = nil
		errors.Walk(merr, func(err error) (stop bool) {
			//t.Logf("%v", err)
			errs = append(errs, err)
			return false
		})
		// NOTE: it is not 2 because multi error itself is also an error ...
		assert.Equal(t, 1+2, len(errs), "WalkFunc is called same times as length of error list")
	})

	t.Run("nil", func(t *testing.T) {
		var err error
		walked := 0
		errors.Walk(err, func(err error) (stop bool) {
			walked++
			return false
		})
		assert.Equal(t, 0, walked)
	})
}
