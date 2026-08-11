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
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockWingmanService struct {
	shouldError bool
	errMessage  string
}

func (m *mockWingmanService) CreateReferral(ctx context.Context, wingmanID, target1ID, target2ID uuid.UUID) (*models.WingmanReferral, error) {
	if m.shouldError {
		return nil, errors.New("mock create error")
	}
	return &models.WingmanReferral{
		ID:        uuid.New(),
		DeepLink:  "qlove://match/123",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}, nil
}

func (m *mockWingmanService) AcceptReferral(ctx context.Context, referralID, acceptingUserID uuid.UUID) (*models.WingmanReferral, error) {
	if m.shouldError {
		return nil, errors.New(m.errMessage)
	}
	return &models.WingmanReferral{
		ID:     referralID,
		Status: "matched",
	}, nil
}

func (m *mockWingmanService) ProcessCommission(ctx context.Context, referralID uuid.UUID) error {
	return nil
}

func TestCreateReferral(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockWingmanService{}
	h := NewWingmanHandler(mockSvc)
	app.Post("/wingmans/referral", h.CreateReferral)

	t.Run("Success", func(t *testing.T) {
		payload := CreateReferralRequest{
			Target1ID: uuid.New(),
			Target2ID: uuid.New(),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/wingmans/referral", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil || resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/wingmans/referral", bytes.NewReader([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		appErr := fiber.New()
		errSvc := &mockWingmanService{shouldError: true}
		hErr := NewWingmanHandler(errSvc)
		appErr.Post("/wingmans/referral", hErr.CreateReferral)

		payload := CreateReferralRequest{
			Target1ID: uuid.New(),
			Target2ID: uuid.New(),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/wingmans/referral", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := appErr.Test(req, -1)
		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}

func TestAcceptReferral(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockWingmanService{}
	h := NewWingmanHandler(mockSvc)
	app.Post("/wingmans/referral/:id/accept", h.AcceptReferral)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/wingmans/referral/"+uuid.New().String()+"/accept", nil)
		req.Header.Set("X-User-ID", uuid.New().String())

		resp, err := app.Test(req, -1)
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/wingmans/referral/invalid-uuid/accept", nil)

		resp, _ := app.Test(req, -1)
		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Service Bad Request", func(t *testing.T) {
		appErr := fiber.New()
		errSvc := &mockWingmanService{shouldError: true, errMessage: "referral link expired"}
		hErr := NewWingmanHandler(errSvc)
		appErr.Post("/wingmans/referral/:id/accept", hErr.AcceptReferral)

		req := httptest.NewRequest("POST", "/wingmans/referral/"+uuid.New().String()+"/accept", nil)
		resp, _ := appErr.Test(req, -1)
		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Service Internal Error", func(t *testing.T) {
		appErr := fiber.New()
		errSvc := &mockWingmanService{shouldError: true, errMessage: "database connection lost"}
		hErr := NewWingmanHandler(errSvc)
		appErr.Post("/wingmans/referral/:id/accept", hErr.AcceptReferral)

		req := httptest.NewRequest("POST", "/wingmans/referral/"+uuid.New().String()+"/accept", nil)
		resp, _ := appErr.Test(req, -1)
		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}
