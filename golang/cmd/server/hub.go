package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var once sync.Once
var hub *Hub

type Hub struct {
	rooms    map[string]*Room
	clients  map[string]*User
	upgrader websocket.Upgrader
	mu       sync.Mutex
}

// NewHub creates a new WebSocket server
func NewHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			upgrader: websocket.Upgrader{},
			rooms:    make(map[string]*Room),
			clients:  make(map[string]*User),
		}

		go func() {
			for range time.Tick(5 * time.Minute) {
				hub.removeEmptyRooms()
			}
		}()
	})
	return hub
}

// handleConnections upgrades HTTP connections to WebSocket and handles new connections
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Fatalf("error upgrading http connection to ws: %v", err)
	}
	defer conn.Close()

	// Parse the form to get the room name
	err = r.ParseForm()
	if err != nil {
		zap.S().Errorf("error parsing form: %v", err)
		return
	}

	roomName := r.Form.Get("room")
	if roomName == "" {
		roomName = "general"
	}

	// Retrieve the user's token from the cookies
	var user *User
	cookie, _ := r.Cookie("token")
	if cookie == nil {
		user = h.createUser(conn)
	} else {
		user = h.getUserByToken(cookie.Value)
		if user == nil {
			user = h.createUser(conn)
		}
	}

	// Set the token cookie, renewing the expiration date
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   user.token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	room := h.joinRoom(user, roomName)

	// Main message handling loop
	for {
		var msg Message
		err := user.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				zap.S().Infof("user %s disconnected", user.name)
				room.users[user] = false
			} else {
				zap.S().Errorf("error reading json from message: %v", err)
			}
			break
		}

		room.broadcastMessage(user.token, msg)
	}
}

func (h *Hub) removeEmptyRooms() {
	h.mu.Lock()
	for roomName, room := range h.rooms {
		if len(room.users) == 0 {
			room.closeChan <- true
			delete(h.rooms, roomName)
			zap.S().Infof("room %s removed for inactivity", roomName)
		}
	}
	h.mu.Unlock()
}
