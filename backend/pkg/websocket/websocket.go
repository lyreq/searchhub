package websocket

import (
	"sync"
)

type Event struct {
	DeviceID  string `json:"device_id,omitempty"`
	IsSub     int    `json:"is_sub,omitempty"`
	Code      string `json:"code,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Payload struct {
	clients    map[string]chan Event
	clientsMux sync.RWMutex
}

func New(clients map[string]chan Event) *Payload {
	return &Payload{
		clients: clients,
	}
}

func (s *Payload) RegisterClient(clientID string) chan Event {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	ch := make(chan Event, 100)
	s.clients[clientID] = ch
	return ch
}

func (s *Payload) UnregisterClient(clientID string) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	if ch, exists := s.clients[clientID]; exists {
		close(ch)
		delete(s.clients, clientID)
	}
}

func (s *Payload) Broadcast(event Event) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	for _, ch := range s.clients {
		select {
		case ch <- event:
		default:
			// Channel is full, skip this client
		}
	}
}

func (s *Payload) BroadcastToClient(clientID string, event Event) {
	s.clientsMux.RLock()
	ch, exists := s.clients[clientID]
	s.clientsMux.RUnlock()

	if !exists {
		return
	}

	select {
	case ch <- event:
	default:
		s.UnregisterClient(clientID)
	}
}

func (s *Payload) FindByToken(token string) bool {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	_, exists := s.clients[token]
	return exists
}
