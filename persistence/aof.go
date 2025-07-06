package persistence

import (
	"bufio" 
	"log"
	"os"
	"sync"

	"github.com/HORUSCRIME/goredis/command" 
	"github.com/HORUSCRIME/goredis/database" 
	"github.com/HORUSCRIME/goredis/resp"     
)

type AOF struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex 
}

func NewAOF(filename string) (*AOF, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &AOF{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (a *AOF) AppendCommand(cmd resp.Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()


	log.Printf("AOF: Appending command: %v", cmd)

	return nil
}

func (a *AOF) Load(db *database.Database, processor *command.Processor) error {


	log.Println("AOF: Loading data from AOF file (placeholder).")
	return nil
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file != nil {
		if err := a.writer.Flush(); err != nil {
			log.Printf("Error flushing AOF writer: %v", err)
		}
		return a.file.Close()
	}
	return nil
}
