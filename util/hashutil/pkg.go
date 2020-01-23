// Package hashutil provides alloc free alternatives for pkg/hash
package hashutil

func HashStringFnv64a(str string) uint64 {
	h := NewInlineFNV64a()
	if _, err := h.WriteString(str); err != nil {
		panic(err)
	}
	return h.Sum64()
}

func HashFnv64a(b []byte) uint64 {
	h := NewInlineFNV64a()
	if _, err := h.Write(b); err != nil {
		panic(err)
	}
	return h.Sum64()
}

// TODO:
// https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
// https://github.com/segmentio/fasthash
