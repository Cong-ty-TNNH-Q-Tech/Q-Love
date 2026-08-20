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
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockVoucherServiceForAdmin struct {
	createErr error
	getErr    error
	deleteErr error
}

func (m *mockVoucherServiceForAdmin) RedeemVoucher(ctx context.Context, userID uuid.UUID, brand string, valueXu int) (*models.Voucher, error) {
	return nil, nil
}
func (m *mockVoucherServiceForAdmin) CreateVoucher(ctx context.Context, req services.CreateVoucherRequest) error {
	return m.createErr
}
func (m *mockVoucherServiceForAdmin) GetVouchers(ctx context.Context, limit, offset int) ([]models.Voucher, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return []models.Voucher{{ID: uuid.New(), Brand: "Test"}}, nil
}
func (m *mockVoucherServiceForAdmin) DeleteVoucher(ctx context.Context, id uuid.UUID) error {
	return m.deleteErr
}

func TestAdminVoucherHandler_CreateVoucher(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherServiceForAdmin{}
	h := NewAdminVoucherHandler(svc)
	app.Post("/admin/v1/vouchers", h.CreateVoucher)

	// Success
	body, _ := json.Marshal(map[string]interface{}{
		"brand": "Highlands", "code": "HL-1", "value_xu": 100, "expires_at": time.Now(),
	})
	req := httptest.NewRequest(http.MethodPost, "/admin/v1/vouchers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Invalid body
	req = httptest.NewRequest(http.MethodPost, "/admin/v1/vouchers", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Missing fields
	body, _ = json.Marshal(map[string]interface{}{"brand": ""})
	req = httptest.NewRequest(http.MethodPost, "/admin/v1/vouchers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Service error
	svc.createErr = errors.New("lỗi db")
	body, _ = json.Marshal(map[string]interface{}{
		"brand": "Highlands", "code": "HL-1", "value_xu": 100, "expires_at": time.Now(),
	})
	req = httptest.NewRequest(http.MethodPost, "/admin/v1/vouchers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAdminVoucherHandler_GetVouchers(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherServiceForAdmin{}
	h := NewAdminVoucherHandler(svc)
	app.Get("/admin/v1/vouchers", h.GetVouchers)

	req := httptest.NewRequest(http.MethodGet, "/admin/v1/vouchers", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	svc.getErr = errors.New("lỗi db")
	req = httptest.NewRequest(http.MethodGet, "/admin/v1/vouchers", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAdminVoucherHandler_DeleteVoucher(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherServiceForAdmin{}
	h := NewAdminVoucherHandler(svc)
	app.Delete("/admin/v1/vouchers/:id", h.DeleteVoucher)

	id := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/admin/v1/vouchers/"+id, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Invalid ID
	req = httptest.NewRequest(http.MethodDelete, "/admin/v1/vouchers/invalid", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Service error
	svc.deleteErr = errors.New("lỗi db")
	req = httptest.NewRequest(http.MethodDelete, "/admin/v1/vouchers/"+id, nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
