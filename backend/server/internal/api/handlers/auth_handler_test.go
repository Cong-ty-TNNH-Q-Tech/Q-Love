// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mocks

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) SendOTP(ctx context.Context, phone string) error {
	args := m.Called(ctx, phone)
	return args.Error(0)
}

func (m *mockAuthService) VerifyOTP(ctx context.Context, phone, otp string) (*models.User, string, string, bool, error) {
	args := m.Called(ctx, phone, otp)
	if args.Get(0) == nil {
		return nil, "", "", false, args.Error(4)
	}
	return args.Get(0).(*models.User), args.String(1), args.String(2), args.Bool(3), args.Error(4)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

func setupAuthTestApp(svc *mockAuthService) *fiber.App {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}
	app := fiber.New()
	handler := NewAuthHandler(svc)
	app.Post("/auth/send-otp", handler.SendOTP)
	app.Post("/auth/verify-otp", handler.VerifyOTP)
	app.Post("/auth/refresh", handler.RefreshToken)
	return app
}

func TestAuthHandler_SendOTP_Success(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("SendOTP", mock.Anything, "0901234567").Return(nil)

	reqBody, _ := json.Marshal(map[string]string{"phone": "0901234567"})
	req := httptest.NewRequest(http.MethodPost, "/auth/send-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_SendOTP_RateLimit(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("SendOTP", mock.Anything, "0901234567").Return(errors.New("rate limit exceeded"))

	reqBody, _ := json.Marshal(map[string]string{"phone": "0901234567"})
	req := httptest.NewRequest(http.MethodPost, "/auth/send-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestAuthHandler_SendOTP_InternalError(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("SendOTP", mock.Anything, "0901234567").Return(errors.New("some unexpected error"))

	reqBody, _ := json.Marshal(map[string]string{"phone": "0901234567"})
	req := httptest.NewRequest(http.MethodPost, "/auth/send-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAuthHandler_SendOTP_InvalidBody(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	req := httptest.NewRequest(http.MethodPost, "/auth/send-otp", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAuthHandler_VerifyOTP_Success(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	user := &models.User{ID: uuid.New(), Phone: "0901234567"}
	svc.On("VerifyOTP", mock.Anything, "0901234567", "123456").Return(user, "access", "refresh", false, nil)

	reqBody, _ := json.Marshal(map[string]string{"phone": "0901234567", "otp": "123456"})
	req := httptest.NewRequest(http.MethodPost, "/auth/verify-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_VerifyOTP_InvalidOTP(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("VerifyOTP", mock.Anything, "0901234567", "123456").Return(nil, "", "", false, errors.New("ERR_INVALID_OTP"))

	reqBody, _ := json.Marshal(map[string]string{"phone": "0901234567", "otp": "123456"})
	req := httptest.NewRequest(http.MethodPost, "/auth/verify-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("RefreshToken", mock.Anything, "old-refresh").Return("new-access", "new-refresh", nil)

	reqBody, _ := json.Marshal(map[string]string{"refresh_token": "old-refresh"})
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_RefreshToken_Invalid(t *testing.T) {
	svc := new(mockAuthService)
	app := setupAuthTestApp(svc)

	svc.On("RefreshToken", mock.Anything, "bad").Return("", "", errors.New("invalid"))

	reqBody, _ := json.Marshal(map[string]string{"refresh_token": "bad"})
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
