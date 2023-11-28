package main

import (
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var once sync.Once
var hub *Hub

type Hub struct {
	rooms    map[string]*Room
	clients  map[net.Addr]*User
	upgrader websocket.Upgrader
	mu       sync.Mutex
}

// NewHub creates a new WebSocket server
func NewHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			upgrader: websocket.Upgrader{},
			rooms:    make(map[string]*Room),
			clients:  make(map[net.Addr]*User),
		}
	})
	return hub
}

// handleConnections upgrades HTTP connections to WebSocket and handles new connections
func (s *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Fatalf("error upgrading http connection to ws: %v", err)
	}

	var roomName string
	err = r.ParseForm()
	if err != nil {
		roomName = "general"
	} else {
		roomName = r.Form.Get("room")
	}

	user := s.createUser(conn)
	s.joinRoom(user, roomName)

	// Main message handling loop
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			zap.S().Errorf("error reading json from message: %v", err)
		}

		user.handleMessage(roomName, msg)
	}
}
