package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
		Type:       "text",
		Content:    "hello",
	}
	
	err := hub.PublishMessage(ctx, uuid.New(), msg)
	if err != nil {
		t.Errorf("Expected nil error for nil redis, got %v", err)
	}
}

func TestHub_ConsumeRedisStream(t *testing.T) {
	logger.InitLogger("test", "")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts, _ := redis.ParseURL("redis://localhost:6379/0")
	rdb := redis.NewClient(opts)
	// Only proceed if redis is up (like in CI)
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis is not available, skipping test")
	}
	// Clear the stream for clean test
	rdb.Del(ctx, RedisChatStream)

	hub := NewHub(rdb)
	
	client := &Client{
		Hub:    hub,
		UserID: uuid.New(),
		Send:   make(chan []byte, 256),
	}
	
	go hub.Run(ctx)
	
	// Register client
	hub.Register <- client
	time.Sleep(50 * time.Millisecond)

	// Publish message targeting this client
	msg := &models.ChatMessage{
		ID:      uuid.New(),
		MatchID: uuid.New(),
		Type:    "text",
		Content: "redis stream test",
	}
	
	err := hub.PublishMessage(ctx, client.UserID, msg)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait to receive message from the stream
	select {
	case received := <-client.Send:
		if len(received) == 0 {
			t.Errorf("Expected message, got empty")
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for message from Redis stream")
	}
	
	// Unregister
	hub.Unregister <- client
	time.Sleep(50 * time.Millisecond)
	
	// Test graceful shutdown
	cancel()
	time.Sleep(100 * time.Millisecond)
}
