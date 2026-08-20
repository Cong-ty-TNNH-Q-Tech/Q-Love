package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type mockShameService struct {
	getActiveShamesFn func(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error)
	throwTomatoFn     func(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error
}

func (m *mockShameService) GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error) {
	return m.getActiveShamesFn(ctx, limit, offset)
}

func (m *mockShameService) ThrowTomato(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error {
	return m.throwTomatoFn(ctx, throwerID, shameID)
}

func TestShameHandler_GetActiveShames(t *testing.T) {
	app := fiber.New()

	mockSvc := &mockShameService{
		getActiveShamesFn: func(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error) {
			return []models.WallOfShameResponse{
				{WallOfShame: models.WallOfShame{ID: uuid.New(), UserID: uuid.New(), Reason: "Test", TomatoesThrown: 10, ExpiresAt: time.Now().Add(1 * time.Hour)}, UserName: "TestUser", AvatarURL: "test.jpg"},
			}, nil
		},
	}
	h := NewShameHandler(mockSvc)
	app.Get("/shames", h.GetActiveShames)

	req := httptest.NewRequest(http.MethodGet, "/shames?limit=10&offset=0", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Test error case
	mockSvc.getActiveShamesFn = func(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error) {
		return nil, errors.New("db error")
	}
	respErr, _ := app.Test(req)
	if respErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", respErr.StatusCode)
	}
}

func TestShameHandler_ThrowTomato(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockShameService{
		throwTomatoFn: func(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error {
			return nil
		},
	}
	h := NewShameHandler(mockSvc)
	app.Post("/shames/:id/tomato", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return c.Next()
	}, h.ThrowTomato)

	shameID := uuid.New().String()
	reqBody := map[string]string{"thrower_id": uuid.New().String()}
	bodyBytes, _ := json.Marshal(reqBody)

	// Success case
	req := httptest.NewRequest(http.MethodPost, "/shames/"+shameID+"/tomato", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Invalid ID
	reqInvalidID := httptest.NewRequest(http.MethodPost, "/shames/invalid-id/tomato", bytes.NewBuffer(bodyBytes))
	reqInvalidID.Header.Set("Content-Type", "application/json")
	respInvalidID, _ := app.Test(reqInvalidID)
	if respInvalidID.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected Bad Request, got %v", respInvalidID.StatusCode)
	}

	// Service error (Insufficient Balance)
	mockSvc.throwTomatoFn = func(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error {
		return services.ErrInsufficientBalance
	}
	respInsuff, _ := app.Test(req)
	if respInsuff.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected Bad Request, got %v", respInsuff.StatusCode)
	}

	// Service error (Other)
	mockSvc.throwTomatoFn = func(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error {
		return errors.New("internal error")
	}
	respErr, _ := app.Test(req)
	if respErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected Internal Server Error, got %v", respErr.StatusCode)
	}
}
