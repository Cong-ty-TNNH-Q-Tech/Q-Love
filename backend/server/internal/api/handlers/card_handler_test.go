// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCardService struct {
	mock.Mock
}

func (m *MockCardService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CardProfile), args.Error(1)
}

func (m *MockCardService) TradeCard(ctx context.Context, collectorID, targetUserID uuid.UUID, tradeType string, quantity int) error {
	args := m.Called(ctx, collectorID, targetUserID, tradeType, quantity)
	return args.Error(0)
}

func setupCardApp(svc *MockCardService) *fiber.App {
	app := fiber.New()
	handler := NewCardHandler(svc)
	
	// Mock locals middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "123e4567-e89b-12d3-a456-426614174000")
		return c.Next()
	})
	
	app.Get("/cards/:user_id", handler.GetCardProfile)
	app.Post("/cards/trade", handler.TradeCard)
	
	return app
}

func TestGetCardProfile(t *testing.T) {
	svc := new(MockCardService)
	app := setupCardApp(svc)
	
	targetID := uuid.New()
	
	// Success
	svc.On("GetProfile", mock.Anything, targetID).Return(&models.CardProfile{
		UserID: targetID,
	}, nil).Once()
	
	req := httptest.NewRequest("GET", "/cards/"+targetID.String(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Invalid UUID
	req2 := httptest.NewRequest("GET", "/cards/invalid-uuid", nil)
	resp2, _ := app.Test(req2)
	assert.Equal(t, 422, resp2.StatusCode)
	
	// Internal Error
	svc.On("GetProfile", mock.Anything, targetID).Return(nil, errors.New("db error")).Once()
	req3 := httptest.NewRequest("GET", "/cards/"+targetID.String(), nil)
	resp3, _ := app.Test(req3)
	assert.Equal(t, 500, resp3.StatusCode)
}

func TestTradeCard(t *testing.T) {
	svc := new(MockCardService)
	app := setupCardApp(svc)
	
	targetID := uuid.New()
	
	// Success
	svc.On("TradeCard", mock.Anything, mock.Anything, targetID, "buy", 1).Return(nil).Once()
	
	body, _ := json.Marshal(TradeCardRequest{
		TargetUserID: targetID.String(),
		Type:         "buy",
		Quantity:     1,
	})
	req := httptest.NewRequest("POST", "/cards/trade", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Circuit Breaker error
	svc.On("TradeCard", mock.Anything, mock.Anything, targetID, "buy", 1).Return(errors.New("circuit breaker active")).Once()
	req2 := httptest.NewRequest("POST", "/cards/trade", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	assert.Equal(t, 503, resp2.StatusCode)
	
	// Level 5 error
	svc.On("TradeCard", mock.Anything, mock.Anything, targetID, "buy", 1).Return(errors.New("need level 5")).Once()
	req3 := httptest.NewRequest("POST", "/cards/trade", bytes.NewReader(body))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := app.Test(req3)
	assert.Equal(t, 403, resp3.StatusCode)
	
	// Bad Request error
	svc.On("TradeCard", mock.Anything, mock.Anything, targetID, "buy", 1).Return(errors.New("other error")).Once()
	req4 := httptest.NewRequest("POST", "/cards/trade", bytes.NewReader(body))
	req4.Header.Set("Content-Type", "application/json")
	resp4, _ := app.Test(req4)
	assert.Equal(t, 400, resp4.StatusCode)
	
	// Invalid Target UUID
	bodyInvalid, _ := json.Marshal(TradeCardRequest{
		TargetUserID: "invalid-uuid",
		Type:         "buy",
		Quantity:     1,
	})
	req5 := httptest.NewRequest("POST", "/cards/trade", bytes.NewReader(bodyInvalid))
	req5.Header.Set("Content-Type", "application/json")
	resp5, _ := app.Test(req5)
	assert.Equal(t, 422, resp5.StatusCode)
}
