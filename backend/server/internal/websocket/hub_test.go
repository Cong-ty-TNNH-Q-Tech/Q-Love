package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
)

func TestHub_Run(t *testing.T) {
	logger.InitLogger("test", "")
	// Create a new hub with nil redis client (which is handled gracefully now)
	hub := NewHub(nil)
	
	// Start the hub in a background goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go hub.Run(ctx)
	
	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Test PublishMessage with nil client
	msg := &models.ChatMessage{
		ID:         uuid.New(),
		SenderID:   uuid.New(),
		Content:    "Hello",
	}
	err := hub.PublishMessage(ctx, uuid.New(), msg)
	if err != nil {
		t.Errorf("Expected nil error when redis client is nil, got %v", err)
	}
	
	// Test graceful shutdown
	cancel()
	time.Sleep(100 * time.Millisecond)
}
