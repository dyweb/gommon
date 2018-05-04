package cli

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestHandler_HandleLog(t *testing.T) {
	// FIXME: the test seems to be broken, it has no output and when run in GoLand it shows terminated instead of pass
	//fmt.Printf("%04d", 2)
}

func Test_FormatNum(t *testing.T) {
	assert := asst.New(t)

	assert.Equal("0010", string(formatNum(10, 4)))
	assert.Equal("0100", string(formatNum(100, 4)))
	assert.Equal("1000", string(formatNum(1000, 4)))
	assert.Equal("0000", string(formatNum(10000, 4)))
}
