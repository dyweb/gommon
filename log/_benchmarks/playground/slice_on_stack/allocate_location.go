package main

import (
	"io"
)

type handle struct {
	w io.Writer

	// value receiver
	f  f
	fp *f

	// pointer receiver
	f2  f2
	f2p *f2

	// interface
	fi fi
}

type f struct {
}

//go:noinline
func (f f) Write(b []byte) {
	b[0] = 'c'
}

type f2 struct {
}

type fi interface {
	Write(b []byte)
}

//go:noinline
func (f *f2) Write(b []byte) {
	b[0] = 'c'
}

func (h *handle) Log(s string) {
	l1 := make([]byte, 0, 100) // constant type, but still escape to heap,
	l1 = append(l1, s...)
	h.w.Write(l1) // parameter to indirect call
}

func (h *handle) Log2(s string) {
	l2 := make([]byte, 0, len(s)) // escape (non-constant size)
	l2 = append(l2, s...)
	h.w.Write(l2)
}

func (h *handle) Log3(s string) {
	l3 := make([]byte, 0, 100) // don't escape to heap
	l3 = append(l3, s...)
	h.f.Write(l3)
}

func (h *handle) Log4(s string) {
	l4 := make([]byte, 0, 100) // don't escape to heap TODO: because receiver is not pointer?
	l4 = append(l4, s...)
	h.fp.Write(l4)
}

func (h *handle) Log5(s string) {
	l5 := make([]byte, 0, 100) // still don't escape to heap ...
	l5 = append(l5, s...)
	h.f2.Write(l5)
}

func (h *handle) Log6(s string) {
	l6 := make([]byte, 0, 100) // still don't escape
	l6 = append(l6, s...)
	h.f2p.Write(l6)
}

func (h *handle) Log7(s string) {
	l7 := make([]byte, 0, 100) // escape ...
	l7 = append(l7, s...)
	h.fi.Write(l7)
}

//go:noinline
func (h *handle) Compute(i int) {
	c1 := make([]byte, 0, 100) // don't escape to heap
	c1[i] = 'c'
}

//go:noinline
func (h *handle) Compute2(i int) {
	c1 := make([]byte, i) // escape to heap, unknown size
	c1[i-1] = 'c'
}

// go build -gcflags "-m -m" .
func main() {
}
