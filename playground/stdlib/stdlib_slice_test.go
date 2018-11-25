package stdlib

import (
	"testing"

	asst "github.com/stretchr/testify/assert"
)

type foo struct {
	bar string
	m   map[string]string
	s   []string
}

// When shallow copy a slice of struct, all fields are copied by value, including map,
// but map itself is reference to underlying structure, so if you update map in old struct, you will see the new value in new struct as well
// If you re-slice and assign new struct, it won't have effect on the copied strut's map because it's still referencing the old data
func TestSlice_ReSliceStructSlice(t *testing.T) {
	assert := asst.New(t)

	s := make([]foo, 0)
	m := map[string]string{"a": "1"}
	ss := []string{"1", "2"}
	s = append(s, foo{"1", m, ss})
	s = append(s, foo{"2", m, ss})

	// shallow copy
	sCopy := make([]foo, len(s))
	copy(sCopy, s)

	// map is a reference to underlying data structure
	s[0].bar = "1+"
	s[0].m["a"] = "1+"
	assert.Equal("1", sCopy[0].bar)
	assert.Equal("1+", sCopy[0].m["a"])
	// slice is also a reference
	s[0].s[0] = "1+"
	s[0].s = append(s[0].s, "3")
	assert.NotEqual(3, len(sCopy[0].s)) // extend the slice in original struct won't have effect on the length of the copy
	assert.Equal("1+", sCopy[0].s[0])

	// re-slice and assign new value won't effect map in the copied slice
	s = s[:0]
	mr := map[string]string{"a": "1r"}
	s = append(s, foo{"1r", mr, ss})
	s = append(s, foo{"2r", mr, ss})
	assert.Equal("1+", sCopy[0].m["a"])
}

func acceptVariadicArgs(nums ...int) int {
	s := 0
	for _, v := range nums {
		s += v
	}
	return s
}

func TestVariadicArgs(t *testing.T) {
	v1 := []int{1, 2, 3}
	v2 := []int{4, 5, 6}
	acceptVariadicArgs(v1...)
	acceptVariadicArgs(v2...)
	//acceptVariadicArgs(v2..., v1...)
}
