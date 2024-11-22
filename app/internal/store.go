package internal

import "sync"

type Store struct {
	kv map[string]string
	mu sync.Mutex
}

func (s *Store) Insert(k string, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.kv[k] = v

}
