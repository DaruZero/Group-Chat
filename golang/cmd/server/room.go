package main

import "go.uber.org/zap"

type Room struct {
	name     string
	users    map[*User]bool
	messages []Message
}

// joinRoom adds a user to a room
func (h *Hub) joinRoom(user *User, roomName string) {
	h.mu.Lock()

	// Create a new room if it doesn't exist
	var room *Room
	var ok bool
	if room, ok = h.rooms[roomName]; !ok {
		h.mu.Unlock()
		room = h.createRoom(roomName)
		h.mu.Lock()
	}

	room.users[user] = true
	h.mu.Unlock()

	zap.S().Infof("user %s joined room %s", user.name, roomName)
}

// createRoom creates a new room. It overwrites an existing room if it has the same name.
func (h *Hub) createRoom(roomName string) *Room {
	zap.S().Infof("creating room %s", roomName)

	h.mu.Lock()
	defer h.mu.Unlock()

	room := &Room{
		name:     roomName,
		users:    make(map[*User]bool),
		messages: make([]Message, 0),
	}
	h.rooms[roomName] = room
	zap.S().Infof("created room %s", roomName)
	return room
}
