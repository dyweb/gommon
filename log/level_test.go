package log2

import (
	"fmt"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	assert := asst.New(t)
	cases := []struct {
		Str string
		Lvl Level
		Num int
	}{
		{"fatal", FatalLevel, 0},
		{"panic", PanicLevel, 1},
		{"error", ErrorLevel, 2},
		{"warn", WarnLevel, 3},
		{"info", InfoLevel, 4},
		{"debug", DebugLevel, 5},
		{"trace", TraceLevel, 6},
	}

	for _, c := range cases {
		assert.Equal(c.Num, int(c.Lvl))
		assert.Equal(c.Str, fmt.Sprint(c.Lvl))
	}
}
