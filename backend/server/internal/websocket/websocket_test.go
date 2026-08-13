package websocket

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	fiberWebsocket "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func TestWebsocket_Integration(t *testing.T) {
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
	app.Get("/ws", fiberWebsocket.New(func(c *fiberWebsocket.Conn) {
		serverClient = &Client{
			Hub:    hub,
			Conn:   c,
			UserID: uuid.New(),
			Send:   make(chan []byte, 256),
		}
		serverClient.Hub.Register <- serverClient
		go serverClient.WritePump()
		serverClient.ReadPump()
	}))

	// Start server on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go app.Listener(ln)
	defer app.Shutdown()

	time.Sleep(100 * time.Millisecond)

	url := fmt.Sprintf("ws://127.0.0.1:%d/ws", port)
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Wait for register
	time.Sleep(100 * time.Millisecond)

	// Test sending a message FROM the server TO the client
	if serverClient != nil {
		serverClient.Send <- []byte("hello client")
	}

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
