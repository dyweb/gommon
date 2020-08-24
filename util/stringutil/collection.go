package stringutil

import "sort"

// collection.go contains set etc.

type Set struct {
	m         map[string]struct{}
	insertion []string
}

func NewSet() *Set {
	return &Set{
		m: make(map[string]struct{}),
	}
}

func (s *Set) Add(str string) {
	_, ok := s.m[str]
	if ok {
		return
	}
	s.m[str] = struct{}{}
	s.insertion = append(s.insertion, str)
}

func (s *Set) AddMulti(ss ...string) {
	for _, str := range ss {
		s.Add(str)
	}
}

// Inserted returns a copy of unique strings in their insertion order.
func (s *Set) Inserted() []string {
	return CopySlice(s.insertion)
}

// Sorted returns sorted unique strings in ascending order. e.g. [a, b, c]
func (s *Set) Sorted() []string {
	cp := CopySlice(s.insertion)
	sort.Strings(cp)
	return cp
}

// TODO: impl
//func (s *Set) SortedDesc() []string  {
//
//}
