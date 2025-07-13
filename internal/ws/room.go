package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Room *Room
	Send chan []byte
}

type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.RWMutex
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Rooms[client.Room.ID]; !ok {
		h.Rooms[client.Room.ID] = &Room{
			ID: client.Room.ID,
			Clients: make(map[*Client]bool),
		}
	}

	room := h.Rooms[client.Room.ID]
	room.mu.Lock()
	room.Clients[client] = true
	room.mu.Unlock()
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.Rooms[client.Room.ID]; ok {
		room.mu.Lock()
		if _, ok := room.Clients[client]; ok {
			delete(room.Clients, client)
			close(client.Send)
		}
		room.mu.Unlock()

		if len(room.Clients) == 0 {
			delete(h.Rooms, client.Room.ID)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte, sender *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.Rooms[roomID]; ok {
		room.mu.RLock()
		clientCount := len(room.Clients)
		
		sentCount := 0
		for client := range room.Clients {
			if client != sender {
				select {
				case client.Send <- message:
					sentCount++
				default:
					log.Printf("Failed to send to client %s, removing", client.ID)
					close(client.Send)
					delete(room.Clients, client)
				}
			}
		}
		room.mu.RUnlock()
	} else {
		log.Printf("Room %s not found for broadcast", roomID)
	}
}

func (h *Hub) GetRoomInfo(roomID string) (int, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.Rooms[roomID]; ok {
		room.mu.RLock()
		count := len(room.Clients)
		room.mu.RUnlock()
		return count, true
	}
	return 0, false
}