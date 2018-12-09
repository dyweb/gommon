package cli

import (
	"fmt"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestFmt_FormatNum(t *testing.T) {
	fmt.Printf("%04d\n", 2)
}

func Test_FormatNum(t *testing.T) {
	assert := asst.New(t)

	assert.Equal("0010", string(formatNum(10, 4)))
	assert.Equal("0100", string(formatNum(100, 4)))
	assert.Equal("1000", string(formatNum(1000, 4)))
	assert.Equal("0000", string(formatNum(10000, 4)))
}
