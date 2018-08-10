package errors

import (
	"fmt"
	"os"
	"testing"

	"github.com/dyweb/gommon/util/testutil"
	asst "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := asst.New(t)
	err := New("don't let me go")
	assert.NotNil(err)
	assert.Equal("don't let me go", fmt.Sprintf("%v", err))
	terr, ok := err.(TracedError)
	assert.True(ok)
	printFrames(terr.ErrorStack().Frames())
	assert.Equal(3, len(terr.ErrorStack().Frames()))
}

func freshErr() error {
	return New("I am a fresh error")
}

func wrappedStdErr() error {
	return Wrap(os.ErrClosed, "can't open closed file")
}

func TestWrap(t *testing.T) {
	assert := asst.New(t)

	n := Wrap(nil, "nothing")
	assert.Nil(n)

	errw := Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal("can't open closed file: file already closed", fmt.Sprintf("%v", errw))
	terr, ok := errw.(TracedError)
	assert.True(ok)
	if testutil.Dump().B() {
		printFrames(terr.ErrorStack().Frames())
	}
	assert.Equal(3, len(terr.ErrorStack().Frames()))

	errw = Wrap(freshErr(), "wrap again")
	terr, ok = errw.(TracedError)
	assert.True(ok)
	if testutil.Dump().B() {
		printFrames(terr.ErrorStack().Frames())
	}
	assert.Equal(4, len(terr.ErrorStack().Frames()))

	errw = Wrap(wrappedStdErr(), "wrap again")
	terr, ok = errw.(TracedError)
	assert.True(ok)
	if testutil.Dump().B() {
		printFrames(terr.ErrorStack().Frames())
	}
	assert.Equal(4, len(terr.ErrorStack().Frames()))
}

// TODO: test Wrapf

func TestCause(t *testing.T) {
	assert := asst.New(t)
	n := Wrap(nil, "nothing")
	assert.Nil(Cause(n))

	errw := Wrap(os.ErrClosed, "can't open closed file")
	assert.Equal(os.ErrClosed, Cause(errw))

	errww := Wrap(errw, "wrap again")
	assert.Equal(os.ErrClosed, Cause(errww))
}

func TestDirectCause(t *testing.T) {
	assert := asst.New(t)

	errw := Wrap(os.ErrClosed, "can't open closed file")
	errww := Wrap(errw, "wrap again")
	assert.Equal(os.ErrClosed, Cause(errww))
	assert.Equal(os.ErrClosed, DirectCause(errw))
	assert.NotEqual(os.ErrClosed, DirectCause(errww))
	assert.Equal("can't open closed file", DirectCause(errww).(Wrapper).Message())
}

func TestWrappedError_Message(t *testing.T) {
	assert := asst.New(t)

	msg := "mewo"
	errw := Wrap(os.ErrClosed, msg)
	assert.Equal(msg, errw.(Wrapper).Message())
}

// TODO: we should write test in errors_test package, especially for examples ...
func ExampleWrap() {
	err := Wrap(os.ErrNotExist, "oops")
	fmt.Println(err)
	// Output:
	// oops: file does not exist
}
