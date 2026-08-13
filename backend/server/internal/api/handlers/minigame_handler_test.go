// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockMinigameService struct{}

func (m *mockMinigameService) InitSteal(ctx context.Context, attackerID uuid.UUID, defenderID uuid.UUID, targetCardID uuid.UUID) (*models.CardSteal, error) {
	return &models.CardSteal{ID: uuid.New()}, nil
}

func (m *mockMinigameService) SubmitStealResult(ctx context.Context, stealID uuid.UUID, attackerID uuid.UUID, isWin bool) error {
	return nil
}

func TestMinigameHandler_InitSteal(t *testing.T) {
	app := fiber.New()
	handler := NewMinigameHandler(&mockMinigameService{})

	app.Post("/init", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.InitSteal(c)
	})

	reqBody := map[string]string{
		"defender_id":    uuid.New().String(),
		"target_card_id": uuid.New().String(),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/init", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected 200 OK, got %v", resp.StatusCode)
	}
}

func TestMinigameHandler_SubmitStealResult(t *testing.T) {
	app := fiber.New()
	handler := NewMinigameHandler(&mockMinigameService{})

	app.Post("/submit", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SubmitStealResult(c)
	})

	reqBody := map[string]interface{}{
		"steal_id": uuid.New().String(),
		"is_win":   true,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected 200 OK, got %v", resp.StatusCode)
	}
}
