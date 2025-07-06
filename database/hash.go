package database

import "sync"

type Hash struct {
	mu sync.RWMutex
	data map[string]string
}

func NewHash() *Hash {
	return &Hash{
		data: make(map[string]string),
	}
}

func (h *Hash) Type() string {
	return "hash"
}

func (h *Hash) HSet(field, value string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, exists := h.data[field]
	h.data[field] = value
	return !exists 
}

func (h *Hash) HGet(field string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	val, ok := h.data[field]
	return val, ok
}

func (h *Hash) HDel(fields ...string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	deletedCount := 0
	for _, field := range fields {
		if _, ok := h.data[field]; ok {
			delete(h.data, field)
			deletedCount++
		}
	}
	return deletedCount
}

func (h *Hash) HLen() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}