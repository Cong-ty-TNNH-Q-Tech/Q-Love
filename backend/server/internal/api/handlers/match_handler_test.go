// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.


package handlers

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockMatchService struct {
	err error
}

func (m *mockMatchService) Unmatch(ctx context.Context, matchID, userID uuid.UUID) error {
	return m.err
}

func TestMatchHandler_Unmatch_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockMatchService{}
	handler := NewMatchHandler(mockService)

	app.Delete("/matches/:match_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.Unmatch(c)
	})

	req := httptest.NewRequest("DELETE", "/matches/"+uuid.NewString(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestMatchHandler_Unmatch_InvalidMatchID(t *testing.T) {
	app := fiber.New()
	mockService := &mockMatchService{}
	handler := NewMatchHandler(mockService)

	app.Delete("/matches/:match_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.Unmatch(c)
	})

	req := httptest.NewRequest("DELETE", "/matches/invalid-uuid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestMatchHandler_Unmatch_Forbidden(t *testing.T) {
	app := fiber.New()
	mockService := &mockMatchService{err: errors.New("forbidden")}
	handler := NewMatchHandler(mockService)

	app.Delete("/matches/:match_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.Unmatch(c)
	})

	req := httptest.NewRequest("DELETE", "/matches/"+uuid.NewString(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

func TestMatchHandler_Unmatch_NotFound(t *testing.T) {
	app := fiber.New()
	mockService := &mockMatchService{err: errors.New("match not found")}
	handler := NewMatchHandler(mockService)

	app.Delete("/matches/:match_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.Unmatch(c)
	})

	req := httptest.NewRequest("DELETE", "/matches/"+uuid.NewString(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestMatchHandler_Unmatch_InternalError(t *testing.T) {
	app := fiber.New()
	mockService := &mockMatchService{err: errors.New("db error")}
	handler := NewMatchHandler(mockService)

	app.Delete("/matches/:match_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.Unmatch(c)
	})

	req := httptest.NewRequest("DELETE", "/matches/"+uuid.NewString(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

