package transaction

import (
	"fmt" 
	"log"
	"sync"
)

type Transaction struct {
	mu        sync.RWMutex 
	commands  []interface{} 
	watching  []string      
	discarded bool
}

func NewTransaction() *Transaction {
	return &Transaction{
		commands: make([]interface{}, 0),
		watching: make([]string, 0),
	}
}

func (t *Transaction) EnqueueCommand(cmd interface{}) {
	t.mu.Lock() 
	defer t.mu.Unlock()
	t.commands = append(t.commands, cmd)
	log.Printf("Transaction: Enqueued command: %v", cmd)
}

func (t *Transaction) WatchKey(key string) {
	t.mu.Lock() 
	defer t.mu.Unlock()
	t.watching = append(t.watching, key)
	log.Printf("Transaction: Watching key: %s", key)
}

func (t *Transaction) Discard() {
	t.mu.Lock() 
	defer t.mu.Unlock()
	t.commands = nil
	t.watching = nil
	t.discarded = true
	log.Println("Transaction: Discarded.")
}

func (t *Transaction) IsDiscarded() bool {
	t.mu.RLock() 
	defer t.mu.RUnlock()
	return t.discarded
}

func (t *Transaction) Execute(db interface{}, processor interface{}) ([]interface{}, error) {
	t.mu.Lock() 
	defer t.mu.Unlock()

	if t.discarded {
		return nil, fmt.Errorf("transaction already discarded")
	}


	log.Printf("Transaction: Executing %d commands (placeholder).", len(t.commands))
	results := make([]interface{}, len(t.commands))
	for i, cmd := range t.commands {

		log.Printf("  Executing: %v", cmd)
		results[i] = "OK (simulated)" 
	}

	t.commands = nil 
	t.watching = nil
	return results, nil
}