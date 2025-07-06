package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HORUSCRIME/goredis/server"
)

func main() {
	address := "0.0.0.0:6379" 

	s := server.NewServer(address) 
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("GoRedis server started on address %s", address) 

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down GoRedis server...")
	s.Stop()
	log.Println("GoRedis server stopped.")
}