package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/HORUSCRIME/goredis/command"
	"github.com/HORUSCRIME/goredis/database"
)

type Server struct {
	//port      string
	listener  net.Listener
	clients   map[*Client]bool
	mu        sync.RWMutex
	db        *database.Database
	processor *command.Processor
	shutdown  chan struct{}

	address string
}

// // NewServer creates a new GoRedis server.
// func NewServer(port string) *Server {
// 	db := database.NewDatabase()
// 	return &Server{
// 		port:      port,
// 		clients:   make(map[*Client]bool),
// 		db:        db,
// 		processor: command.NewProcessor(db), // Initialize command processor with database
// 		shutdown:  make(chan struct{}),
// 	}
// }

// // Start starts the server and listens for incoming connections.
// func (s *Server) Start() error {
// 	addr := fmt.Sprintf(":%s", s.port)
// 	listener, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return fmt.Errorf("failed to listen on %s: %w", addr, err)
// 	}
// 	s.listener = listener
// 	log.Printf("Server listening on %s", addr)

// 	go s.acceptConnections()
// 	return nil
// }

func NewServer(address string) *Server {
	db := database.NewDatabase()
	return &Server{
		address:   address,
		clients:   make(map[*Client]bool),
		db:        db,
		processor: command.NewProcessor(db),
		shutdown:  make(chan struct{}),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}
	s.listener = listener
	log.Printf("Server listening on %s", s.address)

	go s.acceptConnections()
	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				log.Println("Server listener closed, stopping acceptConnections.")
				return
			default:
				log.Printf("Error accepting connection: %v", err)
				continue
			}
		}

		log.Printf("New client connected: %s", conn.RemoteAddr())
		client := NewClient(conn, s.processor)
		s.addClient(client)
		go client.Handle()
	}
}

func (s *Server) Stop() {
	close(s.shutdown)
	if s.listener != nil {
		s.listener.Close()
	}

	s.mu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[*Client]bool)
	s.mu.Unlock()

	log.Println("Server stopped.")
}

func (s *Server) addClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = true
}

func (s *Server) removeClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client)
}
