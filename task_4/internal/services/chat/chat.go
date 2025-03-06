package chat

import "github.com/gorilla/websocket"

type Service struct {
	hub *hub
}

// Create new Service and run hub
func NewService() *Service {
	hub := newHub()
	go hub.run()

	return &Service{hub: hub}
}

func (s *Service) Register(name string, conn *websocket.Conn) error {
	newClient(s.hub, name, conn)
	return nil
}
