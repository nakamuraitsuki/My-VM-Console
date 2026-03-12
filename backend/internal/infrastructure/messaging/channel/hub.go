package channel

import "sync"

type Hub struct {
	mu     *sync.RWMutex
	channels map[string][]chan []byte
}

func NewHub() *Hub {
	return &Hub{
		channels: make(map[string][]chan []byte),
	}
}
