package testutil

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestGenGolden(t *testing.T) {
	// export GOLDEN=true
	// export GEN_GOLDEN=true
	t.Run("generate", func(t *testing.T) {
		RunIf(t, GenGolden())
		WriteFixture(t, "testdata/golden", []byte("I am a file"))
	})
	t.Run("use", func(t *testing.T) {
		assert := asst.New(t)
		SkipIf(t, GenGolden())
		d := ReadFixture(t, "testdata/golden")
		assert.Equal("I am a file", string(d))
	})
}

func TestGenGoldenT(t *testing.T) {
	txt := "gen golden t"
	file := "testdata/goldent"
	t.Run("generate", func(t *testing.T) {
		EnableGenGolden(t)
		RunIf(t, GenGoldenT(t))
		WriteFixture(t, file, []byte(txt))
	})
	t.Run("compare", func(t *testing.T) {
		WriteOrCompare(t, file, []byte(txt))
	})
}
