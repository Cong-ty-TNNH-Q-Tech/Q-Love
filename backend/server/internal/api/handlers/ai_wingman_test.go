// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/ai"
)

type mockAIWingmanService struct {
	getSuggestionsFn func(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error)
}

func (m *mockAIWingmanService) GetSuggestions(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error) {
	return m.getSuggestionsFn(ctx, matchID)
}

func TestAIWingmanHandler_GetSuggestions(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAIWingmanService{
		getSuggestionsFn: func(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error) {
			return []ai.Suggestion{
				{Tone: "Hài hước", Text: "Haha"},
			}, nil
		},
	}
	h := NewAIWingmanHandler(mockSvc)
	app.Get("/matches/:id/wingman-suggestions", h.GetSuggestions)

	matchID := uuid.New().String()

	// Success case
	req := httptest.NewRequest(http.MethodGet, "/matches/"+matchID+"/wingman-suggestions", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Invalid ID
	reqInvalidID := httptest.NewRequest(http.MethodGet, "/matches/invalid-id/wingman-suggestions", nil)
	respInvalidID, _ := app.Test(reqInvalidID)
	if respInvalidID.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected Bad Request, got %v", respInvalidID.StatusCode)
	}

	// Service error
	mockSvc.getSuggestionsFn = func(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error) {
		return nil, errors.New("ai error")
	}
	respErr, _ := app.Test(req)
	if respErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected Internal Server Error, got %v", respErr.StatusCode)
	}
}
