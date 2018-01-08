package stdlib

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

type ptrFoo struct {
}

func TestPtrCompare(t *testing.T) {
	assert := asst.New(t)

	o := ptrFoo{}
	p1 := &o
	p2 := &o
	p3 := &ptrFoo{}
	assert.True(p1 == p2)
	assert.False(p1 == p3)
	assert.False(p2 == p3)
	p3 = p2
	assert.True(p1 == p3)
}
