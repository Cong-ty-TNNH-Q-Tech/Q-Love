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

func TestAdminHandler_Errors(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminService{}
	handler := NewAdminHandler(mockSvc)
	
	app.Put("/admin/ban/:id", handler.BanUser)
	app.Delete("/admin/violations/:id/media", handler.DeleteViolationMedia)
	app.Post("/admin/court/:id/override", handler.OverrideCourtCase)

	// BanUser Invalid UUID
	req1 := httptest.NewRequest("PUT", "/admin/ban/invalid-uuid", nil)
	resp1, _ := app.Test(req1)
	assert.Equal(t, fiber.StatusBadRequest, resp1.StatusCode)

	// DeleteViolationMedia Invalid UUID
	req2 := httptest.NewRequest("DELETE", "/admin/violations/invalid-uuid/media", nil)
	resp2, _ := app.Test(req2)
	assert.Equal(t, fiber.StatusBadRequest, resp2.StatusCode)

	// DeleteViolationMedia Invalid JSON
	req3 := httptest.NewRequest("DELETE", "/admin/violations/"+uuid.New().String()+"/media", bytes.NewReader([]byte(`invalid`)))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := app.Test(req3)
	assert.Equal(t, fiber.StatusBadRequest, resp3.StatusCode)

	// OverrideCourtCase Invalid UUID
	req4 := httptest.NewRequest("POST", "/admin/court/invalid-uuid/override", nil)
	resp4, _ := app.Test(req4)
	assert.Equal(t, fiber.StatusBadRequest, resp4.StatusCode)

	// OverrideCourtCase Invalid JSON
	req5 := httptest.NewRequest("POST", "/admin/court/"+uuid.New().String()+"/override", bytes.NewReader([]byte(`invalid`)))
	req5.Header.Set("Content-Type", "application/json")
	resp5, _ := app.Test(req5)
	assert.Equal(t, fiber.StatusBadRequest, resp5.StatusCode)
}

type mockAdminServiceError struct {
	mockAdminService
}
func (m *mockAdminServiceError) BanUser(ctx context.Context, userID uuid.UUID) error {
	return errors.New("service error")
}
func (m *mockAdminServiceError) OverrideCourtCase(ctx context.Context, courtID uuid.UUID, status string) error {
	return errors.New("service error")
}
func (m *mockAdminServiceError) DeleteViolationMedia(ctx context.Context, violationID uuid.UUID, objectKey string) error {
	return errors.New("service error")
}
func (m *mockAdminServiceError) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	return nil, 0, errors.New("service error")
}

func TestAdminHandler_ServiceErrors(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockAdminServiceError{}
	handler := NewAdminHandler(mockSvc)
	
	app.Get("/admin/violations", handler.GetViolations)
	app.Put("/admin/ban/:id", handler.BanUser)
	app.Delete("/admin/violations/:id/media", handler.DeleteViolationMedia)
	app.Post("/admin/court/:id/override", handler.OverrideCourtCase)

	req0 := httptest.NewRequest("GET", "/admin/violations", nil)
	resp0, _ := app.Test(req0)
	assert.Equal(t, fiber.StatusInternalServerError, resp0.StatusCode)

	req1 := httptest.NewRequest("PUT", "/admin/ban/"+uuid.New().String(), nil)
	resp1, _ := app.Test(req1)
	assert.Equal(t, fiber.StatusInternalServerError, resp1.StatusCode)

	body2, _ := json.Marshal(map[string]string{"object_key": "test.jpg"})
	req2 := httptest.NewRequest("DELETE", "/admin/violations/"+uuid.New().String()+"/media", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	assert.Equal(t, fiber.StatusInternalServerError, resp2.StatusCode)

	body3, _ := json.Marshal(map[string]string{"status": "dismissed"})
	req3 := httptest.NewRequest("POST", "/admin/court/"+uuid.New().String()+"/override", bytes.NewReader(body3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := app.Test(req3)
	assert.Equal(t, fiber.StatusInternalServerError, resp3.StatusCode)
}
