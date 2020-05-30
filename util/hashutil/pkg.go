// Package hashutil provides alloc free alternatives for pkg/hash.
// Currently only fnv64a and fnv32a are supported.
// TODO:
// https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
// https://github.com/segmentio/fasthash
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

func HashFnv32a(b []byte) uint32 {
	h := NewInlineFNV32a()
	if _, err := h.Write(b); err != nil {
		panic(err)
	}
	return h.Sum32()
}
