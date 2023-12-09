package main

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Room struct {
	users     map[*User]bool
	closeChan chan bool
	name      string
	messages  []Message
	mu        sync.Mutex
}

type Message struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// joinRoom adds a user to a room and returns the room.
func (h *Hub) joinRoom(user *User, roomName string) *Room {

	// Create a new room if it doesn't exist
	var room *Room
	var ok bool
	h.mu.Lock()
	if room, ok = h.rooms[roomName]; !ok {
		h.mu.Unlock()
		room = h.createRoom(roomName)
	}

	room.mu.Lock()
	room.users[user] = true

	// Send chat history to the new user
	for _, msg := range room.messages {
		user.conn.WriteJSON(msg)
	}
	room.mu.Unlock()

	zap.S().Infof("user %s joined room %s", user.name, roomName)

	return room
}

// createRoom creates a new room. It overwrites an existing room if it has the same name.
func (h *Hub) createRoom(roomName string) *Room {
	zap.S().Infof("creating room %s", roomName)

	room := &Room{
		name:     roomName,
		users:    make(map[*User]bool),
		messages: make([]Message, 0),
	}

	h.mu.Lock()
	h.rooms[roomName] = room
	h.mu.Unlock()

	zap.S().Infof("created room %s", roomName)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				room.removeInactiveUsers()
			case <-room.closeChan:
				ticker.Stop()
				return
			}
		}
	}()

	return room
}

func (r *Room) broadcastMessage(sender string, message Message) {
	message.Sender = sender

	r.mu.Lock()
	r.messages = append(r.messages, message)
	r.mu.Unlock()

	for user := range r.users {
		if user.token != sender {
			user.conn.WriteJSON(message)
		}
	}
}

func (r *Room) removeInactiveUsers() {
	now := time.Now()
	r.mu.Lock()
	for user, active := range r.users {
		if !active && now.Sub(user.lastActive) > 1*time.Hour {
			delete(r.users, user)
			zap.S().Infof("user %s removed from room %s for inactivity", user.name, r.name)
		}
	}
	r.mu.Unlock()
}
