package mem

import (
	"github.com/dyweb/gommon/log2/_examples/uselib/storage"
)

type Store struct {
	d map[string]string
}

func (s *Store) Get(k string) (string, error) {
	return s.d[k], nil
}

func (s *Store) Set(k string, v string) {
	s.d[k] = v
}

func NewMemStorage() *Store {
	return &Store{
		d: make(map[string]string, 5),
	}
}

func init() {
	storage.Register("mem", func() storage.Driver {
		return NewMemStorage()
	})
}
