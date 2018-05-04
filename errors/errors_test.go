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

// TODO: we should write test in errors_test package, especially for examples ...
func ExampleWrap() {
	err := Wrap(os.ErrNotExist, "oops")
	fmt.Println(err)
	// Output:
	// oops: file does not exist
}
