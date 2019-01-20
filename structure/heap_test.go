package structure

import "testing"

// TODO: assertion on edge cases ... and why GoLand is not having hint? due to dep mode enabled?
func TestIntHeap_Insert(t *testing.T) {
	h := IntHeap{}
	h.Insert(2)
	h.Insert(1)
	t.Log(h)
	t.Log(h.Poll())
	t.Log(h)

	h2 := IntHeap{}
	unsorted := []int{3, 1, 4, 5, 6, -1, 2, 8, 9, 1}
	for i := 0; i < len(unsorted); i++ {
		h2.Insert(unsorted[i])
	}
	var sorted []int
	for h2.Len() != 0 {
		sorted = append(sorted, h2.Poll())
	}
	t.Log(sorted)
}
