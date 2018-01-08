package stdlib

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

type ptrFoo struct {
	// NOTE: this attribute is needed, because otherwise in go1.7 and 1.8 emtpy struct seems to be optimized by compiler
	// to use same value, and the pointer to and empty struct would be same
	// p1 := &ptrFoo{}
	// p2 := &ptrFoo{}
	// p1 == p2 // true in 1.7 and 1.8, false in 1.9 and tip
	// https://github.com/dyweb/gommon/issues/36
	foo string
}

func TestPtrCompare(t *testing.T) {
	assert := asst.New(t)

	o := ptrFoo{foo: "o"}
	p1 := &o
	p2 := &o
	p3 := &ptrFoo{foo: "p3"}
	assert.True(p1 == p2)
	assert.False(p1 == p3)
	assert.False(p2 == p3)
	p3 = p2
	assert.True(p1 == p3)
}
