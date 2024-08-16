package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	clients    map[string]*websocket.Conn
	clientsMux sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *WebSocketService) AddClient(userID string, conn *websocket.Conn) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	s.clients[userID] = conn
}

func (s *WebSocketService) RemoveClient(userID string) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	delete(s.clients, userID)
}

func (s *WebSocketService) SendUpdate(userID string, message []byte) {
	s.clientsMux.RLock()
	conn, exists := s.clients[userID]
	s.clientsMux.RUnlock()

	if exists {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
		}
	}
}
