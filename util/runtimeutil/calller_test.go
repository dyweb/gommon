package runtimeutil

import (
	asst "github.com/stretchr/testify/assert"
	"testing"
)

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
