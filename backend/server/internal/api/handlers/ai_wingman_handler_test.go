// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAIWingmanService struct{
	err error
}

func (m *mockAIWingmanService) SuggestReplies(ctx context.Context, matchID uuid.UUID) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []string{"Suggest 1", "Suggest 2", "Suggest 3"}, nil
}
func (m *mockAIWingmanService) MaskPII(text string) string { return text }

func TestAIWingmanHandler_SuggestReplies(t *testing.T) {
	app := fiber.New()
	handler := NewAIWingmanHandler(&mockAIWingmanService{})

	app.Post("/suggest", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SuggestReplies(c)
	})

	matchID := uuid.New().String()
	body := map[string]string{"match_id": matchID}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/suggest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result struct {
		Replies []string `json:"replies"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Len(t, result.Replies, 3)
	assert.Equal(t, "Suggest 1", result.Replies[0])
}

func TestAIWingmanHandler_SuggestReplies_InvalidJSON(t *testing.T) {
	app := fiber.New()
	handler := NewAIWingmanHandler(&mockAIWingmanService{})
	app.Post("/suggest", handler.SuggestReplies)

	req := httptest.NewRequest(http.MethodPost, "/suggest", bytes.NewReader([]byte(`invalid`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAIWingmanHandler_SuggestReplies_InvalidUUID(t *testing.T) {
	app := fiber.New()
	handler := NewAIWingmanHandler(&mockAIWingmanService{})
	app.Post("/suggest", handler.SuggestReplies)

	body := map[string]string{"match_id": "invalid-uuid"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/suggest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAIWingmanHandler_SuggestReplies_ServiceError(t *testing.T) {
	app := fiber.New()
	handler := NewAIWingmanHandler(&mockAIWingmanService{err: assert.AnError})
	app.Post("/suggest", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SuggestReplies(c)
	})

	body := map[string]string{"match_id": uuid.New().String()}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/suggest", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
