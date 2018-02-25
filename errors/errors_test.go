package errors

import (
	"os"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := asst.New(t)
	err := New("don't let me go")
	assert.NotNil(err)
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

	errw := Wrap(os.ErrClosed, "can't open closed file")
	terr, ok := errw.(TracedError)
	assert.True(ok)
	printFrames(terr.ErrorStack().Frames())
	assert.Equal(3, len(terr.ErrorStack().Frames()))

	errw = Wrap(freshErr(), "wrap again")
	terr, ok = errw.(TracedError)
	assert.True(ok)
	printFrames(terr.ErrorStack().Frames())
	assert.Equal(4, len(terr.ErrorStack().Frames()))

	errw = Wrap(wrappedStdErr(), "wrap again")
	terr, ok = errw.(TracedError)
	assert.True(ok)
	printFrames(terr.ErrorStack().Frames())
	assert.Equal(4, len(terr.ErrorStack().Frames()))
}
