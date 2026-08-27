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
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDatingContractService struct {
	mock.Mock
}

func (m *mockDatingContractService) CreateContract(ctx context.Context, initiatorID, targetUserID uuid.UUID, amount float64, apptTime time.Time) (*models.DatingContract, error) {
	args := m.Called(ctx, initiatorID, targetUserID, amount, apptTime)
	var p *models.DatingContract
	if args.Get(0) != nil {
		p = args.Get(0).(*models.DatingContract)
	}
	return p, args.Error(1)
}

func (m *mockDatingContractService) AcceptContract(ctx context.Context, contractID, targetUserID uuid.UUID) (*models.DatingContract, error) {
	args := m.Called(ctx, contractID, targetUserID)
	var p *models.DatingContract
	if args.Get(0) != nil {
		p = args.Get(0).(*models.DatingContract)
	}
	return p, args.Error(1)
}

func (m *mockDatingContractService) CancelContract(ctx context.Context, contractID, cancellerID uuid.UUID, reason string) error {
	return m.Called(ctx, contractID, cancellerID, reason).Error(0)
}

func (m *mockDatingContractService) ScanContract(ctx context.Context, contractID uuid.UUID, qrToken string) error {
	return m.Called(ctx, contractID, qrToken).Error(0)
}

func TestDatingContractHandler_CreateContract(t *testing.T) {
	app := fiber.New()
	mockService := new(mockDatingContractService)
	handler := NewDatingContractHandler(mockService)

	app.Post("/contract", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		return handler.CreateContract(c)
	})

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	targetID := uuid.New()
	apptTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	t.Run("Success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("CreateContract", mock.Anything, userID, targetID, 500.0, mock.Anything).Return(&models.DatingContract{ID: uuid.New()}, nil)

		reqBody := map[string]interface{}{
			"target_user_id":   targetID.String(),
			"deposit_amount":   500.0,
			"appointment_time": apptTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("CreateContract", mock.Anything, userID, targetID, 500.0, mock.Anything).Return(nil, assert.AnError)

		reqBody := map[string]interface{}{
			"target_user_id":   targetID.String(),
			"deposit_amount":   500.0,
			"appointment_time": apptTime.Format(time.RFC3339),
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDatingContractHandler_AcceptContract(t *testing.T) {
	app := fiber.New()
	mockService := new(mockDatingContractService)
	handler := NewDatingContractHandler(mockService)

	app.Post("/contract/:contract_id/accept", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		return handler.AcceptContract(c)
	})

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	contractID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("AcceptContract", mock.Anything, contractID, userID).Return(&models.DatingContract{ID: contractID}, nil)

		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/accept", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("AcceptContract", mock.Anything, contractID, userID).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/accept", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDatingContractHandler_CancelContract(t *testing.T) {
	app := fiber.New()
	mockService := new(mockDatingContractService)
	handler := NewDatingContractHandler(mockService)

	app.Post("/contract/:contract_id/cancel", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		return handler.CancelContract(c)
	})

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	contractID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("CancelContract", mock.Anything, contractID, userID, "Changed mind").Return(nil)

		reqBody := map[string]interface{}{
			"reason": "Changed mind",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/cancel", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("CancelContract", mock.Anything, contractID, userID, "Changed mind").Return(assert.AnError)

		reqBody := map[string]interface{}{
			"reason": "Changed mind",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/cancel", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDatingContractHandler_ScanContract(t *testing.T) {
	app := fiber.New()
	mockService := new(mockDatingContractService)
	handler := NewDatingContractHandler(mockService)

	app.Post("/contract/:contract_id/scan", handler.ScanContract)

	contractID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("ScanContract", mock.Anything, contractID, "some_qr_token").Return(nil)

		reqBody := map[string]interface{}{
			"qr_token": "some_qr_token",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/scan", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("ScanContract", mock.Anything, contractID, "some_qr_token").Return(assert.AnError)

		reqBody := map[string]interface{}{
			"qr_token": "some_qr_token",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/contract/"+contractID.String()+"/scan", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
