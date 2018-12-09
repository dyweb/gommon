package main

import "testing"

// go test -run none -bench . -benchtime 3s -benchmem -memprofile p.out

func BenchmarkVariadic(b *testing.B) {
	s := 0
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s += fVariadic(1, 2, 3)
		}
	})
}

func BenchmarkSlice(b *testing.B) {
	s := 0
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s += fSlice([]int{1, 2, 3})
		}
	})
}

func BenchmarkVariadicStruct(b *testing.B) {
	s := 0
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s += fVariadicSt(st{x: 1}, st{x: 2})
		}
	})
}

func BenchmarkSliceStruct(b *testing.B) {
	s := 0
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s += fSliceSt([]st{
				{x: 1},
				{x: 2},
			})
		}
	})
}
