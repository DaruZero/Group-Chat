package main

import (
	"github.com/DaruZero/group-chat/golang/internal/tools"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type User struct {
	conn  *websocket.Conn
	name  string
	token string
}

type Message struct {
	sender  string `json:"sender"`
	content string `json:"content"`
}

// createUser creates a new user and adds it to the hub's clients list.
func (h *Hub) createUser(conn *websocket.Conn) *User {
	zap.S().Info("creating new user")
	token := tools.GenerateToken()
	h.mu.Lock()

	if _, ok := h.clients[token]; ok {
		zap.S().Warn("user with token already exists, retrying with new token")
		h.mu.Unlock()
		return h.createUser(conn)
	}

	user := &User{
		conn:  conn,
		name:  namesgenerator.GetRandomName(0),
		token: token,
	}
	h.clients[token] = user
	h.mu.Unlock()
	zap.S().Infof("created user %s", user.name)
	return user
}

// getUserByToken retrieves a user from the hub's clients list by its token.
func (h *Hub) getUserByToken(token string) *User {
	h.mu.Lock()
	defer h.mu.Unlock()

	if user, ok := h.clients[token]; ok {
		return user
	} else {
		zap.S().Warnf("user with token %s not found", token)
		return nil
	}
}

// handleMessage processes an incoming message
func (u *User) handleMessage(roomName string, msg Message) {
	// TODO: handle message
}
