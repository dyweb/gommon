package stdlib

import (
	"testing"
	"io/ioutil"
)

// Both conversion has bytes alloc
// I was thinking one of them should be free, but string is immutable in Go while []byte is mutable
// Came across it both influxdb and uber/zap has custom fnv function to avoid the []byte(string) alloc
// - https://github.com/uber-go/zap/blob/3216c4c73e8ebbafff73c456dd27143f4a7c1f94/zapcore/sampler.go#L51
// - https://github.com/influxdata/influxdb/blob/master/models/inline_fnv.go
// There are blog talking about it
// - https://syslog.ravelin.com/byte-vs-string-in-go-d645b67ca7ff
// TODO: it seems where I put the variable also matters, stack would be used for small strings https://stackoverflow.com/questions/38554773/when-golang-does-allocation-for-string-to-byte-conversion
//
//go test -bench=. -benchmem
//BenchmarkStdString2Bytes-8      50000000                27.8 ns/op            16 B/op          1 allocs/op
//BenchmarkStdBytes2String-8      50000000                28.4 ns/op            16 B/op          1 allocs/op
//BenchmarkStdStringCopy-8        2000000000               0.52 ns/op            0 B/op          0 allocs/op
//PASS
//ok      github.com/dyweb/gommon/playground/stdlib       3.978s


func BenchmarkStdString2Bytes(b *testing.B) {
	s := "I am a string"
	var bs []byte
	for n := 0; n < b.N; n++ {
		bs = []byte(s)
	}
	bs = bs[:1]
	//ioutil.Discard.Write(bs)
}

func BenchmarkStdBytes2String(b *testing.B) {
	bs := []byte("I am a string") // as []byte
	var s string
	for n := 0; n < b.N; n++ {
		s = string(bs)
	}
	s = s[:1]
}

// string is immutable and reference to same underlying bytes when copy http://www.tapirgames.com/blog/golang-string
func BenchmarkStdStringCopy(b *testing.B) {
	s := "I am a string"
	var s2 string
	for n := 0; n < b.N; n++ {
		s2 = s
	}
	ioutil.Discard.Write([]byte(s2))
}