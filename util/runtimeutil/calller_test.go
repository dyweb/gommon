package runtimeutil

import (
	"runtime"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

type Foo struct {
}

func (f Foo) Bar() runtime.Frame {
	return GetCallerFrame(0)
}

type FooPtr struct {
}

func (f *FooPtr) Bar() runtime.Frame {
	return GetCallerFrame(0)
}

func TestGetCallerFrame(t *testing.T) {
	t.Run("value receiver", func(t *testing.T) {
		assert := asst.New(t)
		foo := Foo{}
		fm := foo.Bar()
		assert.Equal("github.com/dyweb/gommon/util/runtimeutil.Foo.Bar", fm.Function)
	})
	t.Run("pointer receiver", func(t *testing.T) {
		assert := asst.New(t)
		foo := FooPtr{}
		fm := foo.Bar()
		assert.Equal("github.com/dyweb/gommon/util/runtimeutil.(*FooPtr).Bar", fm.Function)
	})
}

func TestSplitPackageFunc(t *testing.T) {
	assert := asst.New(t)

	cases := []struct {
		Function     string
		Package      string
		RealFunction string
	}{
		{"github.com/dyweb/gommon/log2/_examples/uselib/service.(*Auth).Check", "github.com/dyweb/gommon/log2/_examples/uselib/service", "(*Auth).Check"},
		{"github.com/dyweb/gommon/log2.TestNewIdentityFromCaller", "github.com/dyweb/gommon/log2", "TestNewIdentityFromCaller"},
	}

	for _, c := range cases {
		p, f := SplitPackageFunc(c.Function)
		assert.Equal(c.Package, p)
		assert.Equal(c.RealFunction, f)
	}
}

func TestSplitStructMethod(t *testing.T) {
	assert := asst.New(t)

	cases := []struct {
		Function string
		Struct   string
		Method   string
	}{
		{"Foo.Bar", "Foo", "Bar"},
		{"(*FooPtr).Bar", "FooPtr", "Bar"},
	}

	for _, c := range cases {
		st, f := SplitStructMethod(c.Function)
		assert.Equal(c.Struct, st)
		assert.Equal(c.Method, f)
	}
}
