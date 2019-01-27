package errors_test

import (
	stderrors "errors"
	"os"
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/stretchr/testify/assert"
)

func TestCause(t *testing.T) {
	n := errors.Wrap(nil, "nothing")
	assert.Nil(t, errors.Cause(n))

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal(t, os.ErrClosed, errors.Cause(errw))

	errww := errors.Wrap(errw, "wrap again")
	assert.Equal(t, os.ErrClosed, errors.Cause(errww))

	stderr := stderrors.New("std")
	assert.Equal(t, stderr, errors.Cause(stderr))
}

func TestDirectCause(t *testing.T) {
	assert.Nil(t, errors.DirectCause(nil))

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	errww := errors.Wrap(errw, "wrap again")
	assert.Equal(t, os.ErrClosed, errors.Cause(errww))
	assert.Equal(t, os.ErrClosed, errors.Cause(errw))
	assert.NotEqual(t, os.ErrClosed, errors.DirectCause(errww))
	assert.Equal(t, "can't open closed file", errors.DirectCause(errww).(errors.Messenger).Message())

	stderr := stderrors.New("std")
	assert.Equal(t, stderr, errors.DirectCause(stderr))
}
