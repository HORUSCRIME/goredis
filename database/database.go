package database

import (
	"log"
	"sync"
	"time"
)

type Database struct {
	mu   sync.RWMutex
	data map[string]Value
	ttl  map[string]time.Time 
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]Value),
		ttl:  make(map[string]time.Time),
	}
}

func (db *Database) Get(key string) (Value, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if expiry, ok := db.ttl[key]; ok && time.Now().After(expiry) {
		log.Printf("Key '%s' expired.", key)
		db.mu.RUnlock() 
		db.Delete(key)
		db.mu.RLock()
		return nil, false
	}

	val, ok := db.data[key]
	return val, ok
}

func (db *Database) Set(key string, val Value, ttl time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = val
	if ttl > 0 {
		db.ttl[key] = time.Now().Add(ttl)
		log.Printf("Set key '%s' with TTL: %v", key, ttl)
	} else {
		delete(db.ttl, key) 
	}
}

func (db *Database) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, exists := db.data[key]
	if exists {
		delete(db.data, key)
		delete(db.ttl, key)
		log.Printf("Key '%s' deleted.", key)
	}
	return exists
}

func (db *Database) Exists(key string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if expiry, ok := db.ttl[key]; ok && time.Now().After(expiry) {
		db.mu.RUnlock()
		db.Delete(key)
		db.mu.RLock()
		return false
	}

	_, ok := db.data[key]
	return ok
}

func (db *Database) Type(key string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.data[key]
	if !ok {
		return "none"
	}
	return val.Type()
}