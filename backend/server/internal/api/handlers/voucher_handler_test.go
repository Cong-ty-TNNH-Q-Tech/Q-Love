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
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockVoucherService struct {
	vouchers []models.Voucher
	redeemV  *models.Voucher
	err      error
}

func (m *mockVoucherService) RedeemVoucher(ctx context.Context, userID uuid.UUID, brand string, valueXu int) (*models.Voucher, error) {
	return m.redeemV, m.err
}
func (m *mockVoucherService) CreateVoucher(ctx context.Context, req services.CreateVoucherRequest) error {
	return m.err
}
func (m *mockVoucherService) GetVouchers(ctx context.Context, limit, offset int) ([]models.Voucher, error) {
	return m.vouchers, m.err
}
func (m *mockVoucherService) DeleteVoucher(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func TestVoucherHandler_RedeemVoucher(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{
		redeemV: &models.Voucher{Brand: "Highlands", ValueXu: 100},
	}
	handler := NewVoucherHandler(svc)
	
	// Add mock locals middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		return c.Next()
	})
	app.Post("/redeem", handler.RedeemVoucher)

	body, _ := json.Marshal(map[string]interface{}{
		"brand":    "Highlands",
		"value_xu": 100,
	})
	req := httptest.NewRequest("POST", "/redeem", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestVoucherHandler_GetAvailableVouchers(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{
		vouchers: []models.Voucher{
			{Brand: "Highlands", ValueXu: 100, Status: "available"},
			{Brand: "Highlands", ValueXu: 100, Status: "available"},
			{Brand: "CGV", ValueXu: 50, Status: "available"},
		},
	}
	handler := NewVoucherHandler(svc)
	app.Get("/vouchers", handler.GetAvailableVouchers)

	req := httptest.NewRequest("GET", "/vouchers", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminVoucherHandler_CreateVoucher(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{}
	handler := NewAdminVoucherHandler(svc)
	app.Post("/vouchers", handler.CreateVoucher)

	body, _ := json.Marshal(map[string]interface{}{
		"brand":      "Highlands",
		"code":       "ABC",
		"value_xu":   100,
		"expires_at": time.Now(),
	})
	req := httptest.NewRequest("POST", "/vouchers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestAdminVoucherHandler_GetVouchers(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{}
	handler := NewAdminVoucherHandler(svc)
	app.Get("/vouchers", handler.GetVouchers)

	req := httptest.NewRequest("GET", "/vouchers", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminVoucherHandler_DeleteVoucher(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{}
	handler := NewAdminVoucherHandler(svc)
	app.Delete("/vouchers/:id", handler.DeleteVoucher)

	req := httptest.NewRequest("DELETE", "/vouchers/"+uuid.New().String(), nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestVoucherHandler_RedeemVoucher_Errors(t *testing.T) {
	app := fiber.New()
	svc := &mockVoucherService{err: errors.New("fail")}
	handler := NewVoucherHandler(svc)
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", uuid.New().String()); return c.Next() })
	app.Post("/redeem", handler.RedeemVoucher)

	body, _ := json.Marshal(map[string]interface{}{
		"brand":    "Highlands",
		"value_xu": 100,
	})
	req := httptest.NewRequest("POST", "/redeem", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
