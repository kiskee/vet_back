package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type ClientRole string

const (
	RoleClient ClientRole = "client"
	RoleVet    ClientRole = "vet"
)

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	byUserID   map[string]*Client
	byVetID    map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[*Client]bool),
		byUserID: make(map[string]*Client),
		byVetID:  make(map[string]*Client),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c] = true
	switch c.role {
	case RoleClient:
		h.byUserID[c.userID] = c
	case RoleVet:
		h.byVetID[c.userID] = c
	}
	log.Printf("ws client registered: role=%s user_id=%s", c.role, c.userID)
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		switch c.role {
		case RoleClient:
			delete(h.byUserID, c.userID)
		case RoleVet:
			delete(h.byVetID, c.userID)
		}
		close(c.send)
		log.Printf("ws client unregistered: role=%s user_id=%s", c.role, c.userID)
	}
}

func (h *Hub) SendToClient(userID string, event Event) {
	h.mu.RLock()
	client, ok := h.byUserID[userID]
	h.mu.RUnlock()

	if ok {
		client.send <- mustMarshal(event)
	}
}

func (h *Hub) SendToVet(vetID string, event Event) {
	h.mu.RLock()
	client, ok := h.byVetID[vetID]
	h.mu.RUnlock()

	if ok {
		client.send <- mustMarshal(event)
	}
}

func (h *Hub) BroadcastToRequest(requestID string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	event.RequestID = requestID
	data := mustMarshal(event)

	for client := range h.clients {
		if client.isSubscribed(requestID) {
			select {
			case client.send <- data:
			default:
				log.Printf("ws client send buffer full: user_id=%s", client.userID)
			}
		}
	}
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("ws marshal error: %v", err)
		return nil
	}
	return data
}
