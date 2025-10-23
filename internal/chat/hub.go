package chat

import (
	"fmt"
	"sync"
)

type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) AddClient(userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = client
	fmt.Println("✅ User connected:", userID)
	fmt.Print("total online users:", len(h.clients), "\n")
}

func (h *Hub) RemoveClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
	fmt.Println("❌ User disconnected:", userID)
	fmt.Print("total online users:", len(h.clients), "\n")
}

func (h *Hub) SendToUser(receiverID string, msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[receiverID]; ok {
		client.Send(msg)
	}
}

func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	_, ok := h.clients[userID]
	h.mu.RUnlock()
	return ok
}
