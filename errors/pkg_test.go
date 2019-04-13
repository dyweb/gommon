package errors_test

import (
	"testing"

	"github.com/dyweb/gommon/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := errors.Errorf("this is a %dd string %s", 3, "!")
	assert.Equal(t, "this is a 3d string !", err.Error())
}

func TestIgnore(t *testing.T) {
	errors.Ignore(errors.Errorf("this is a %d error", 1))
}
