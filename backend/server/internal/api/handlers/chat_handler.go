// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.
package handlers

import (
	"context"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	customWs "github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService services.ChatService
	hub         *customWs.Hub
}

func NewChatHandler(chatService services.ChatService, hub *customWs.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
	}
}

// Upgrade upgrades the HTTP connection to a WebSocket connection
func (h *ChatHandler) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// WSHandler handles the actual websocket connection
func (h *ChatHandler) WSHandler(c *websocket.Conn) {
	// user_id is set by JWT middleware typically, but here we can extract from locals or query for simplicity
	userIDStr := c.Query("user_id") // Simplified for this example. In reality, get from Context / JWT
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Close()
		return
	}

	client := &customWs.Client{
		Hub:    h.hub,
		Conn:   c,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}

	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.WritePump()
	client.ReadPump()
}

// SendMessage API allows sending a message via HTTP which then gets broadcasted via Redis and WS
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	// Extract from JWT in real app, we mock it via headers/query for now
	userIDStr := c.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		MatchID string `json:"match_id"`
		TargetID string `json:"target_id"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	matchID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match_id"})
	}
	
	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid target_id"})
	}

	// 1. Save to DB
	msg, err := h.chatService.SaveMessage(c.Context(), userID, matchID, req.Type, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save message"})
	}

	// 2. Publish to Redis Stream for Real-time delivery
	_ = h.hub.PublishMessage(context.Background(), targetID, msg)

	return c.Status(fiber.StatusCreated).JSON(msg)
}

// GetMessages API returns history of chat
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	matchIDStr := c.Params("match_id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match_id"})
	}

	// Optional before cursor
	beforeStr := c.Query("before")
	var before *time.Time
	if beforeStr != "" {
		parsed, err := time.Parse(time.RFC3339, beforeStr)
		if err == nil {
			before = &parsed
		}
	}

	messages, err := h.chatService.GetMessages(c.Context(), matchID, 50, before)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get messages"})
	}

	return c.JSON(messages)
}

