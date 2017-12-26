package log2

import (
	"testing"
	"fmt"

	asst "github.com/stretchr/testify/assert"
)

type Foo struct {
}

func (f *Foo) LoggerIdentity(justCallMe func() *Identity) *Identity {
	return justCallMe()
}

var _ = NewIdentityFromCaller(0)

func TestNewIdentityFromCaller(t *testing.T) {
	NewIdentityFromCaller(0)
	NewIdentityFromCaller(1)
	NewIdentityFromCaller(2)
}

func TestNewStructLogger(t *testing.T) {
	assert := asst.New(t)
	// FIXME: should pass real package logger
	foo := &Foo{}
	id := foo.LoggerIdentity(NewIdentityFromCaller2)
	assert.Equal("Foo", id.Struct)
	assert.Equal(MagicStructLoggerMethod, id.Function)
	assert.Equal("/home/at15/workspace/src/github.com/dyweb/gommon/log2/identity_test.go:14",
		fmt.Sprintf("%s:%d", id.File, id.Line))
	assert.Equal("struct", id.Type.String())
}
