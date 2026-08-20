// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAdminService struct {
	GetViolationsFn        func(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error)
	BanUserFn              func(ctx context.Context, userID uuid.UUID) error
	DeleteViolationMediaFn func(ctx context.Context, violationID uuid.UUID, objectKey string) error
	OverrideCourtCaseFn    func(ctx context.Context, caseID uuid.UUID, status string) error
}

func (m *mockAdminService) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	if m.GetViolationsFn != nil {
		return m.GetViolationsFn(ctx, page, limit)
	}
	return nil, 0, nil
}
func (m *mockAdminService) BanUser(ctx context.Context, userID uuid.UUID) error {
	if m.BanUserFn != nil {
		return m.BanUserFn(ctx, userID)
	}
	return nil
}
func (m *mockAdminService) DeleteViolationMedia(ctx context.Context, violationID uuid.UUID, objectKey string) error {
	if m.DeleteViolationMediaFn != nil {
		return m.DeleteViolationMediaFn(ctx, violationID, objectKey)
	}
	return nil
}
func (m *mockAdminService) OverrideCourtCase(ctx context.Context, caseID uuid.UUID, status string) error {
	if m.OverrideCourtCaseFn != nil {
		return m.OverrideCourtCaseFn(ctx, caseID, status)
	}
	return nil
}

func TestAdminHandler_GetViolations(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminService{}
	handler := NewAdminHandler(mockSvc)
	app.Get("/admin/violations", handler.GetViolations)

	req := httptest.NewRequest("GET", "/admin/violations?page=1&limit=10", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminHandler_BanUser(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminService{}
	handler := NewAdminHandler(mockSvc)
	app.Post("/admin/users/:id/ban", handler.BanUser)

	req := httptest.NewRequest("POST", "/admin/users/"+uuid.New().String()+"/ban", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminHandler_DeleteViolationMedia(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminService{}
	handler := NewAdminHandler(mockSvc)
	app.Delete("/admin/violations/:id/media", handler.DeleteViolationMedia)

	body, _ := json.Marshal(map[string]string{"object_key": "test.jpg"})
	req := httptest.NewRequest("DELETE", "/admin/violations/"+uuid.New().String()+"/media", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminHandler_OverrideCourtCase(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminService{}
	handler := NewAdminHandler(mockSvc)
	app.Post("/admin/court/:id/override", handler.OverrideCourtCase)

	body, _ := json.Marshal(map[string]string{"status": "dismissed"})
	req := httptest.NewRequest("POST", "/admin/court/"+uuid.New().String()+"/override", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
