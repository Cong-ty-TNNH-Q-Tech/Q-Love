package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	customWs "github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockChatService struct {
	saveMessageFunc func(ctx context.Context, senderID, matchID uuid.UUID, msgType, content string) (*models.ChatMessage, error)
	getMessagesFunc func(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error)
}

func (m *mockChatService) SaveMessage(ctx context.Context, senderID, matchID uuid.UUID, msgType, content string) (*models.ChatMessage, error) {
	if m.saveMessageFunc != nil {
		return m.saveMessageFunc(ctx, senderID, matchID, msgType, content)
	}
	return nil, nil
}

func (m *mockChatService) GetMessages(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	if m.getMessagesFunc != nil {
		return m.getMessagesFunc(ctx, matchID, limit, before)
	}
	return nil, nil
}

func TestChatHandler_Upgrade(t *testing.T) {
	app := fiber.New()
	hub := customWs.NewHub(nil)
	handler := NewChatHandler(&mockChatService{}, hub)
	
	app.Get("/ws", handler.Upgrade)

	req := httptest.NewRequest("GET", "/ws", nil)
	// Normal HTTP request should fail upgrade
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUpgradeRequired {
		t.Errorf("expected 426 Upgrade Required, got %d", resp.StatusCode)
	}
}

func TestChatHandler_SendMessage(t *testing.T) {
	app := fiber.New()
	hub := customWs.NewHub(nil)
	
	mockSvc := &mockChatService{
		saveMessageFunc: func(ctx context.Context, senderID, matchID uuid.UUID, msgType, content string) (*models.ChatMessage, error) {
			return &models.ChatMessage{
				ID:      uuid.New(),
				MatchID: matchID,
				Content: content,
			}, nil
		},
	}
	
	handler := NewChatHandler(mockSvc, hub)
	app.Post("/send", handler.SendMessage)

	t.Run("ValidRequest", func(t *testing.T) {
		body := map[string]interface{}{
			"match_id":  uuid.New().String(),
			"target_id": uuid.New().String(),
			"type":      "text",
			"content":   "hello",
		}
		b, _ := json.Marshal(body)
		
		req := httptest.NewRequest("POST", "/send", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", uuid.New().String())
		
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusCreated {
			t.Errorf("expected 201 Created, got %d", resp.StatusCode)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/send", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/send", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", uuid.New().String())
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})
	
	t.Run("InvalidMatchID", func(t *testing.T) {
		body := map[string]interface{}{
			"match_id":  "invalid",
			"target_id": uuid.New().String(),
			"type":      "text",
			"content":   "hello",
		}
		b, _ := json.Marshal(body)
		
		req := httptest.NewRequest("POST", "/send", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", uuid.New().String())
		
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidTargetID", func(t *testing.T) {
		body := map[string]interface{}{
			"match_id":  uuid.New().String(),
			"target_id": "invalid",
			"type":      "text",
			"content":   "hello",
		}
		b, _ := json.Marshal(body)
		
		req := httptest.NewRequest("POST", "/send", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", uuid.New().String())
		
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})
}

func TestChatHandler_GetMessages(t *testing.T) {
	app := fiber.New()
	hub := customWs.NewHub(nil)
	
	mockSvc := &mockChatService{
		getMessagesFunc: func(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
			return []models.ChatMessage{}, nil
		},
	}
	
	handler := NewChatHandler(mockSvc, hub)
	app.Get("/messages/:match_id", handler.GetMessages)

	t.Run("ValidRequest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/messages/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidMatchID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/messages/invalid", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
		}
	})
	
	t.Run("ValidRequestWithBefore", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/messages/"+uuid.New().String()+"?before=2023-01-01T00:00:00Z", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})
}
