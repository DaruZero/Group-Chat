package main

import "go.uber.org/zap"

type Room struct {
	name     string
	users    map[*User]bool
	messages []Message
}

// joinRoom adds a user to a room
func (s *Hub) joinRoom(user *User, roomName string) {
	s.mu.Lock()

	// Create a new room if it doesn't exist
	var room *Room
	var ok bool
	if room, ok = s.rooms[roomName]; !ok {
		s.mu.Unlock()
		room = s.createRoom(roomName)
		s.mu.Lock()
	}

	room.users[user] = true
	s.mu.Unlock()

	zap.S().Infof("user %s joined room %s", user.name, roomName)
}

// createRoom creates a new room
func (s *Hub) createRoom(roomName string) *Room {
	zap.S().Infof("creating room %s", roomName)

	s.mu.Lock()
	defer s.mu.Unlock()

	if room, ok := s.rooms[roomName]; !ok {
		room = &Room{
			name:     roomName,
			users:    make(map[*User]bool),
			messages: make([]Message, 0),
		}
		s.rooms[roomName] = room
		zap.S().Infof("created room %s", roomName)
		return room
	} else {
		zap.S().Warnf("room %s already exists", roomName)
		return room
	}
}
