package websocket

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	fiberWebsocket "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api/handlers"
)

func TestWebsocket_Integration(t *testing.T) {
	logger.InitLogger("test", "")
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	hub := NewHub(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberWebsocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	var serverClient *Client
	
	// Create dummy chat handler just for WSHandler
	chatHandler := handlers.NewChatHandler(nil, hub)
	
	// We inject user_id in the URL
	testUserID := uuid.New()
	
	app.Get("/ws", fiberWebsocket.New(chatHandler.WSHandler))

	// Start server on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go app.Listener(ln)
	defer app.Shutdown()

	time.Sleep(100 * time.Millisecond)

	url := fmt.Sprintf("ws://127.0.0.1:%d/ws?user_id=%s", port, testUserID.String())
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Wait for register
	time.Sleep(100 * time.Millisecond)

	// We can't access serverClient easily anymore, but we can verify echo or ping
	// In ChatHandler WSHandler, it doesn't echo. It only reads and discards (or handles).
	// Actually we should test publishing a message via the hub.
	msg := []byte(`{"type":"text","content":"hello client"}`)
	hub.PublishMessage(context.Background(), testUserID, string(msg)) // Nil redis client, it will just not use redis stream and rely on local hub routing?
	// Wait, without Redis, Hub.PublishMessage does NOTHING in our mock if we pass nil redis client, or rather it just returns nil.
	// Oh! `PublishMessage` in `hub.go` only pushes to Redis. It does NOT route locally directly. 
	// The routing happens via `consumeRedisStream`.
	// Since we mock it in `websocket_test.go` with `nil` redis, `PublishMessage` won't send it to the client!
	
	// Let's just directly send to the client's channel to test Read/Write Pump.
	hub.mu.RLock()
	for client := range hub.clients[testUserID] {
		client.Send <- msg
	}
	hub.mu.RUnlock()

	// Read message ON the client
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("Failed to read from client: %v", err)
	}
	if msgType != websocket.TextMessage {
		t.Errorf("Expected text message, got %v", msgType)
	}
	if string(msg) != "hello client" {
		t.Errorf("Expected 'hello client', got %s", string(msg))
	}

	// Test sending a message FROM the client TO the server
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"ping": "pong"}`))
	if err != nil {
		t.Errorf("Failed to write message: %v", err)
	}
	
	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	// Close client normally
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Errorf("Failed to write close message: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
}
