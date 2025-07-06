package pubsub

import (
	"log"
	"sync"
)

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan []byte]bool 
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]map[chan []byte]bool),
	}
}

func (ps *PubSub) Subscribe(topic string, clientChan chan []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.subscribers[topic]; !ok {
		ps.subscribers[topic] = make(map[chan []byte]bool)
	}
	ps.subscribers[topic][clientChan] = true
	log.Printf("Subscribed client to topic: %s", topic)
}

func (ps *PubSub) Unsubscribe(topic string, clientChan chan []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if chans, ok := ps.subscribers[topic]; ok {
		delete(chans, clientChan)
		if len(chans) == 0 {
			delete(ps.subscribers, topic)
		}
		log.Printf("Unsubscribed client from topic: %s", topic)
	}
}

func (ps *PubSub) Publish(topic string, message []byte) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	publishedCount := 0
	if chans, ok := ps.subscribers[topic]; ok {
		for clientChan := range chans {
			select {
			case clientChan <- message:
				publishedCount++
			default:
				log.Printf("Skipping message for full channel on topic '%s'", topic)
			}
		}
	}
	log.Printf("Published message to topic '%s', %d subscribers notified.", topic, publishedCount)
	return publishedCount
}