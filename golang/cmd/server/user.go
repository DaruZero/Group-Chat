package main

import (
	"time"

	"github.com/DaruZero/group-chat/golang/internal/tools"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type User struct {
	lastActive time.Time
	conn       *websocket.Conn
	name       string
	token      string
}

// createUser creates a new user and adds it to the hub's clients list.
func (h *Hub) createUser() *User {
	zap.S().Info("creating new user")
	token := tools.GenerateToken()
	h.mu.Lock()

	if _, ok := h.clients[token]; ok {
		zap.S().Warn("user with token already exists, retrying with new token")
		h.mu.Unlock()
		return h.createUser()
	}

	user := &User{
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
