package stringutil_test

import (
	"testing"

	"github.com/dyweb/gommon/util/stringutil"
	"github.com/stretchr/testify/assert"
)

// TODO: fuzz test
func TestSnakeToCamel(t *testing.T) {
	cases := []struct {
		s string
		c string
	}{
		{"snake", "Snake"},
		{"snake_", "Snake"},
		{"snake_case", "SnakeCase"},
		{"snake_case_case", "SnakeCaseCase"},
		{"snake__case", "SnakeCase"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.c, stringutil.SnakeToCamel(tc.s))
	}
}
