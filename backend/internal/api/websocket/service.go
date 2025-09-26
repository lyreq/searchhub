package websocket

import (
	ws "lytemp/pkg/websocket"
)

type Repository interface {
}

type Service struct {
	service *ws.Payload
	repo    Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		service: ws.New(make(map[string]chan ws.Event)),
		repo:    repo,
	}
}

func (s *Service) RegisterClient(clientID string) chan ws.Event {
	return s.service.RegisterClient(clientID)
}

func (s *Service) RegisterBoardQRClient(clientID string) chan ws.Event {
	return s.service.RegisterClient(clientID)
}

func (s *Service) UnregisterClient(clientID string) {
	s.service.UnregisterClient(clientID)
}

func (s *Service) UnregisterBoardQRClient(clientID string) {
	s.service.UnregisterClient(clientID)
}

func (s *Service) Inject(service *ws.Payload) {
	s.service = service
}
func (s *Service) BroadcastToClient(clientID string, event ws.Event) {
	s.service.BroadcastToClient(clientID, event)
}
