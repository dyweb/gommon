// Package hashutil provides alloc free alternatives for pkg/hash
package hashutil // import "github.com/dyweb/gommon/util/hashutil"

func HashStringFnv64a(str string) uint64 {
	h := NewInlineFNV64a()
	h.WriteString(str)
	return h.Sum64()
}

// TODO:
// https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
// https://github.com/segmentio/fasthash
