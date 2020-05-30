package hashutil

const (
	offset32 = 2166136261
	offset64 = 14695981039346656037
	prime32  = 16777619
	prime64  = 1099511628211
)

// InlineFNV64a is a alloc-free version of https://golang.org/pkg/hash/fnv/
// copied from Xephon-K, which is copied from influxdb/models https://github.com/influxdata/influxdb/blob/master/models/inline_fnv.go
type InlineFNV64a uint64

type InlineFNV32a uint32

// NewInlineFNV64a returns a new instance of InlineFNV64a.
func NewInlineFNV64a() InlineFNV64a {
	return offset64
}

func NewInlineFNV32a() InlineFNV32a {
	return offset32
}

// Write adds data to the running hash.
func (s *InlineFNV64a) Write(data []byte) (int, error) {
	hash := uint64(*s)
	for _, c := range data {
		hash ^= uint64(c)
		hash *= prime64
	}
	*s = InlineFNV64a(hash)
	return len(data), nil
}

func (s *InlineFNV32a) Write(data []byte) (int, error) {
	hash := uint32(*s)
	for _, c := range data {
		hash ^= uint32(c)
		hash *= prime32
	}
	*s = InlineFNV32a(hash)
	return len(data), nil
}

// WriteString avoids a []byte(str) conversion BUT yield different result when string contains non ASCII characters
func (s *InlineFNV64a) WriteString(str string) (int, error) {
	hash := uint64(*s)
	for _, c := range str {
		hash ^= uint64(c)
		hash *= prime64
	}
	*s = InlineFNV64a(hash)
	return len(str), nil
}

// Sum64 returns the uint64 of the current resulting hash.
func (s *InlineFNV64a) Sum64() uint64 {
	return uint64(*s)
}

func (s *InlineFNV32a) Sum32() uint32 {
	return uint32(*s)
}
