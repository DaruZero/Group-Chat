package main

import (
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type User struct {
	conn *websocket.Conn
	name string
}

type Message struct {
	sender  string `json:"sender"`
	content string `json:"content"`
}

// createUser creates a new user
func (h *Hub) createUser(conn *websocket.Conn) *User {
	zap.S().Info("creating new user")
	var user *User
	if _, ok := h.clients[conn.RemoteAddr()]; !ok {
		user = &User{
			conn: conn,
			name: namesgenerator.GetRandomName(0),
		}
		h.clients[conn.RemoteAddr()] = user
		zap.S().Infof("created user %s", user.name)
	} else {
		zap.S().Warn("user already exists")
		user = h.clients[conn.RemoteAddr()]
	}
	return user
}

// handleMessage processes an incoming message
func (u *User) handleMessage(roomName string, msg Message) {
	// TODO: handle message
}
