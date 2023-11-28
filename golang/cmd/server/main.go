package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DaruZero/group-chat/golang/internal/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.New("LOG_LEVEL")
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(log)

	zap.S().Info("welcome to the go websocket chat server!")

	// Create the WebSocket hub
	hub := NewHub()

	// Create an HTTP server that listens on the specified port
	srv := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/ws", hub.HandleConnection)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine so it doesn't block.
	go func() {
		zap.S().Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("server error: %v", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	zap.S().Info("shutting down gracefully")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatalf("server forced to shutdown: %v", err)
	}

	zap.S().Info("server exiting")
}
