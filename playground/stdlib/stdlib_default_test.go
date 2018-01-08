package stdlib

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

type fooStr struct {
	str string
}

func TestDefault_StructString(t *testing.T) {
	assert := asst.New(t)

	f := &fooStr{}
	assert.True(f.str == "")
}