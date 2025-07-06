package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/HORUSCRIME/goredis/command"
	"github.com/HORUSCRIME/goredis/resp"
)

type Client struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	processor *command.Processor 
	closed    bool
	mu        sync.Mutex 
}

func NewClient(conn net.Conn, processor *command.Processor) *Client {
	return &Client{
		conn:      conn,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
		processor: processor,
	}
}

func (c *Client) Handle() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Client handler panic: %v", r)
		}
		c.Close()
	}()

	for {
		if c.closed {
			return
		}

		value, err := resp.Decode(c.reader)
		if err != nil {
			if err.Error() == "EOF" { 
				log.Printf("Client disconnected: %s", c.conn.RemoteAddr())
				return
			}
			log.Printf("Error decoding RESP command from %s: %v", c.conn.RemoteAddr(), err)
			c.WriteError(fmt.Sprintf("ERR invalid command: %v", err))
			continue
		}

		response := c.processor.Process(value)

		c.Write(response)
	}
}

func (c *Client) Write(value resp.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	err := resp.Encode(c.writer, value)
	if err != nil {
		log.Printf("Error encoding RESP response to %s: %v", c.conn.RemoteAddr(), err)
		c.Close()
		return
	}
	if err := c.writer.Flush(); err != nil {
		log.Printf("Error flushing writer to %s: %v", c.conn.RemoteAddr(), err)
		c.Close()
	}
}

func (c *Client) WriteError(msg string) {
	c.Write(resp.NewError(msg))
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.conn.Close()
		c.closed = true
		log.Printf("Client connection closed: %s", c.conn.RemoteAddr())
	}
}
