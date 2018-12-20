package errors_test

import (
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
}
