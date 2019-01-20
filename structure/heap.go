package structure

// https://golang.org/src/container/heap/heap.go
// TODO: it seems the type is called Interface .... heap.Interface and heap.Heap ...

// IntHeap is simple min int heap
type IntHeap struct {
	data []int
}

func (h *IntHeap) Len() int {
	return len(h.data)
}

func (h *IntHeap) Insert(x int) {
	h.data = append(h.data, x)
	// bubble up
	h.up(len(h.data) - 1)
}

// Peak returns the root without of removing it
func (h *IntHeap) Peak() int {
	// TODO: need ok to return invalid value?
	if len(h.data) == 0 {
		return -1
	}
	return h.data[0]
}

// Poll return and remove current root
func (h *IntHeap) Poll() int {
	if len(h.data) == 0 {
		return -1
	}

	t := h.data[0]
	if len(h.data) == 1 {
		h.data = h.data[:0]
		return t
	}
	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1] // -1
	h.down(0)
	return t
}

func (h *IntHeap) up(i0 int) {
	i := i0
	for {
		p := (i - 1) / 2
		if h.data[p] <= h.data[i] {
			break
		}
		// swap
		h.data[p], h.data[i] = h.data[i], h.data[p]
		// keep going up
		i = p
	}
}

func (h *IntHeap) down(i0 int) {
	i := i0
	for {
		iv := h.data[i]
		l := 2*i + 1
		if l >= len(h.data)-1 {
			break
		}
		lv := h.data[l]
		r := l + 1
		rv := h.data[r]
		// need to pick between left and right, the smaller one should come up
		if lv <= rv && iv > lv {
			h.data[i] = lv
			h.data[l] = iv
			i = l
			continue
		}
		if rv <= lv && iv > rv {
			h.data[i] = rv
			h.data[r] = iv
			i = r
			continue
		}
		// iv < rv && iv < lv
		break
	}
}
