// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
)

type MockIAPService struct {
	ProcessFunc func(ctx context.Context, event services.RevenueCatEvent) error
}

func (m *MockIAPService) ProcessRevenueCatWebhook(ctx context.Context, event services.RevenueCatEvent) error {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(ctx, event)
	}
	return nil
}

func TestWebhookHandler_HandleRevenueCat(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		RevenueCatWebhookSecret: "test-secret",
	}

	mockSvc := &MockIAPService{
		ProcessFunc: func(ctx context.Context, event services.RevenueCatEvent) error {
			if event.Type == "INVALID" {
				return services.ErrInvalidWebhookPayload
			}
			return nil
		},
	}
	handler := NewWebhookHandler(cfg, mockSvc)
	app.Post("/webhook", handler.HandleRevenueCat)

	// Test 1: Unauthorized
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", resp.StatusCode)
	}

	// Test 2: Valid request
	req = httptest.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte(`{"event":{"type":"NON_RENEWING_PURCHASE"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-secret")
	resp, _ = app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Test 3: Invalid payload triggering 400
	req = httptest.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte(`{"event":{"type":"INVALID"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "test-secret")
	resp, _ = app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}
}
