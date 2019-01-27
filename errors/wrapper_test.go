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

func TestWrappedError_Message(t *testing.T) {
	msg := "mewo"
	errw := errors.Wrap(os.ErrClosed, msg)
	assert.Equal(t, msg, errw.(errors.Messenger).Message())
}



// stubs
func freshErr() error {
	return errors.New("I am a fresh error")
}

func wrappedStdErr() error {
	return errors.Wrap(os.ErrClosed, "can't open closed file")
}
