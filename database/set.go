package database

import "sync"

type Set struct {
	mu   sync.RWMutex
	data map[string]struct{} 
}

func NewSet() *Set {
	return &Set{
		data: make(map[string]struct{}),
	}
}

func (s *Set) Type() string {
	return "set"
}

func (s *Set) SAdd(members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	addedCount := 0
	for _, member := range members {
		if _, ok := s.data[member]; !ok {
			s.data[member] = struct{}{}
			addedCount++
		}
	}
	return addedCount
}

func (s *Set) SRem(members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removedCount := 0
	for _, member := range members {
		if _, ok := s.data[member]; ok {
			delete(s.data, member)
			removedCount++
		}
	}
	return removedCount
}

func (s *Set) SIsMember(member string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[member]
	return ok
}

func (s *Set) SCard() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}