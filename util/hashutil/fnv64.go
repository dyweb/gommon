package hashutil

const (
	prime64  = 1099511628211
	offset64 = 14695981039346656037
)

// InlineFNV64a is a alloc-free version of https://golang.org/pkg/hash/fnv/
// copied from Xephon-K, which is copied from influxdb/models https://github.com/influxdata/influxdb/blob/master/models/inline_fnv.go
type InlineFNV64a uint64

// NewInlineFNV64a returns a new instance of InlineFNV64a.
func NewInlineFNV64a() InlineFNV64a {
	return offset64
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
