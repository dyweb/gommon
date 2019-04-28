// +build ignore

package structure_test

import (
	"sort"
	"testing"

	"github.com/dyweb/gommon/structure"
	"github.com/stretchr/testify/assert"
)

// TODO: assertion on edge cases ... and why GoLand is not having hint? due to dep mode enabled?
func TestIntHeap_Insert(t *testing.T) {
	h := structure.IntHeap{}
	h.Insert(2)
	h.Insert(1)
	t.Log(h)
	t.Log(h.Poll())
	t.Log(h)

	h2 := structure.IntHeap{}
	unsorted := []int{3, 1, 4, 5, 6, -1, 2, 8, 9, 1}
	for i := 0; i < len(unsorted); i++ {
		h2.Insert(unsorted[i])
	}
	var sorted []int
	for h2.Len() != 0 {
		sorted = append(sorted, h2.Poll())
	}
	sort.Ints(unsorted)
	assert.Equal(t, unsorted, sorted)
}
