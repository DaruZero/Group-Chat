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
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Fatalf("error upgrading http connection to ws: %v", err)
	}

	err = r.ParseForm()
	if err != nil {
		zap.S().Errorf("error parsing form: %v", err)
		return
	}

	roomName := r.Form.Get("room")
	if roomName == "" {
		roomName = "general"
	}

	user := h.createUser(conn)
	h.joinRoom(user, roomName)

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
