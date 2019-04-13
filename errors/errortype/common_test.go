package errortype_test

import (
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/dyweb/gommon/errors/errortype"
	"github.com/stretchr/testify/assert"
)

func TestNewNotFound(t *testing.T) {
	err := errortype.NewNotFound("so")
	assert.True(t, errortype.IsNotFound(err))

	err = errors.New("I am just a plain error")
	assert.False(t, errortype.IsNotFound(err))
}

func TestNewAlreadyExists(t *testing.T) {
	err := errortype.NewAlreadyExists("tennis")
	assert.True(t, errortype.IsAlreadyExists(err))

	err = errors.New("I am just a plain error")
	assert.False(t, errortype.IsAlreadyExists(err))
}

func TestNewNotImplemented(t *testing.T) {
	err := errortype.NewNotImplemented("cli")
	assert.True(t, errortype.IsNotImplemented(err))

	err = errors.New("I am just a plain error")
	assert.False(t, errortype.IsNotImplemented(err))
}
