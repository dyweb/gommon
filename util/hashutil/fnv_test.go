package hashutil

import (
	"fmt"
	"hash/fnv"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestNewInlineFNV64a(t *testing.T) {
	assert := asst.New(t)

	h := NewInlineFNV64a()
	h.WriteString("github.com/dyweb/gommon/abc.go:23")
	r1 := h.Sum64()
	h2 := NewInlineFNV64a()
	h2.WriteString("github.com/dyweb/gommon/abc.go:23")
	r2 := h2.Sum64()
	fmt.Println(r1)
	assert.Equal(r1, r2)
}

func TestInlineFNV32a_Write(t *testing.T) {
	data := []string{
		"8KBF520",
		"BJUB",
		"AMD YES!",
	}
	for _, d := range data {
		b := []byte(d)
		stdfnv := fnv.New32a()
		myfnv := NewInlineFNV32a()
		stdfnv.Write(b)
		myfnv.Write(b)
		if stdfnv.Sum32() != myfnv.Sum32() {
			t.Errorf("Mismatch for %s std %d my %d", d, stdfnv.Sum32(), myfnv.Sum32())
		}
	}
}

// for ascii, Write and WriteString has same result, for non-ascii, NO
func TestInlineFNV64a_Write(t *testing.T) {
	assert := asst.New(t)

	// ascii
	s := "github.com/dyweb/gommon/abc.go:23"
	h1, h2 := NewInlineFNV64a(), NewInlineFNV64a()
	h1.Write([]byte(s))
	h2.WriteString(s)
	assert.Equal(h1.Sum64(), h2.Sum64())

	// non-ascii
	cn := "你好世界"
	h3, h4 := NewInlineFNV64a(), NewInlineFNV64a()
	h3.Write([]byte(cn))
	h4.WriteString(cn)
	assert.NotEqual(h3.Sum64(), h4.Sum64())
}

func TestInlineFNV64a_WriteString(t *testing.T) {
	ss := []string{
		"github.com/dyweb/gommon/abc.go:23",
		"github.com/dyweb/gommon/abc.go:24",
		"github.com/dyweb/gommon/abc/de.go:123",
	}
	for _, s := range ss {
		h := NewInlineFNV64a()
		h.WriteString(s)
		fmt.Println(h.Sum64())
	}
}

// TODO: benchmark InlineFNV64a with the hash/fnv
