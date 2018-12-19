package errors_test

import (
	"fmt"
	"os"
	"testing"

	asst "github.com/stretchr/testify/assert"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/util/testutil"
)

func TestNew(t *testing.T) {
	assert := asst.New(t)
	err := errors.New("don't let me go")
	assert.NotNil(err)
	assert.Equal("don't let me go", fmt.Sprintf("%v", err))
	terr, ok := err.(errors.Tracer)
	assert.True(ok)
	errors.PrintFrames(terr.Stack().Frames())
	assert.Equal(3, len(terr.Stack().Frames()))
}

func TestWrap(t *testing.T) {
	assert := asst.New(t)

	n := errors.Wrap(nil, "nothing")
	assert.Nil(n)

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal("can't open closed file: file already closed", fmt.Sprintf("%v", errw))
	terr, ok := errw.(errors.Tracer)
	assert.True(ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(3, len(terr.Stack().Frames()))

	errw = errors.Wrap(freshErr(), "wrap again")
	terr, ok = errw.(errors.Tracer)
	assert.True(ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(4, len(terr.Stack().Frames()))

	errw = errors.Wrap(wrappedStdErr(), "wrap again")
	terr, ok = errw.(errors.Tracer)
	assert.True(ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(4, len(terr.Stack().Frames()))
}

func TestWrapf(t *testing.T) {
	// TODO: need to ensure Wrap and Wrapf are same ...
	assert := asst.New(t)

	// wrap nil return nil
	n := errors.Wrapf(nil, "nothing %d", 2)
	assert.Nil(n)

	// wrap standard error attach stack
	errw := errors.Wrapf(os.ErrClosed, "can't open closed file %s", "gommon.yml")
	assert.Equal("can't open closed file gommon.yml: file already closed", fmt.Sprintf("%v", errw))
	terr, ok := errw.(errors.Tracer)
	assert.True(ok)
	if testutil.Dump().B() {
		errors.PrintFrames(terr.Stack().Frames())
	}
	assert.Equal(3, len(terr.Stack().Frames()))

	// wrap fresh error reuse stack
	ferr := freshErr()
	errw = errors.Wrapf(ferr, "wrap again %d", 2)
	assert.Equal(ferr.(errors.Tracer).Stack(), errw.(errors.Tracer).Stack())

	// wrap wrapped error reuse stack
	errww := errors.Wrapf(errw, "wrap again %d", 3)
	assert.Equal(errw.(errors.Tracer).Stack(), errww.(errors.Tracer).Stack())
}

func TestCause(t *testing.T) {
	assert := asst.New(t)
	n := errors.Wrap(nil, "nothing")
	assert.Nil(errors.Cause(n))

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal(os.ErrClosed, errors.Cause(errw))

	errww := errors.Wrap(errw, "wrap again")
	assert.Equal(os.ErrClosed, errors.Cause(errww))
}

func TestDirectCause(t *testing.T) {
	assert := asst.New(t)

	errw := errors.Wrap(os.ErrClosed, "can't open closed file")
	errww := errors.Wrap(errw, "wrap again")
	assert.Equal(os.ErrClosed, errors.Cause(errww))
	assert.Equal(os.ErrClosed, errors.Cause(errw))
	assert.NotEqual(os.ErrClosed, errors.DirectCause(errww))
	assert.Equal("can't open closed file", errors.DirectCause(errww).(errors.Messenger).Message())
}

func TestWrappedError_Message(t *testing.T) {
	assert := asst.New(t)

	msg := "mewo"
	errw := errors.Wrap(os.ErrClosed, msg)
	assert.Equal(msg, errw.(errors.Messenger).Message())
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
