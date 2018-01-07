// Package hashutil provides alloc free alternatives for pkg/hash
package hashutil

func HashStringFnv64a(str string) uint64 {
	h := NewInlineFNV64a()
	h.WriteString(str)
	return h.Sum64()
}
