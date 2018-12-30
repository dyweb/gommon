// +build ignore

package log

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	assert := asst.New(t)
	// strict
	for _, l := range AllLevels {
		s := l.String()
		l2, err := ParseLevel(s, true)
		assert.Nil(err)
		assert.Equal(l, l2)
	}
	// not strict
	levels := []string{"FA", "Pa", "er", "Warn", "infooo", "debugggg", "Tracer"}
	for i := 0; i < len(AllLevels); i++ {
		l, err := ParseLevel(levels[i], false)
		assert.Nil(err)
		assert.Equal(AllLevels[i], l)
	}
	// invalid
	_, err := ParseLevel("haha", false)
	assert.NotNil(err)
}
