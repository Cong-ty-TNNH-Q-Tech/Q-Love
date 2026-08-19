// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

func setupVibeApp() (*fiber.App, *VibeHandler) {
	// Mock time so that CheckUnlockTime returns true
	services.TimeNow = func() time.Time {
		return time.Date(2026, 1, 1, 23, 0, 0, 0, time.UTC)
	}

	app := fiber.New()
	service := services.NewSpotifyService()
	handler := NewVibeHandler(service)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})

	app.Get("/vibe/status", handler.Status)
	app.Get("/vibe/current-track", handler.CurrentTrack)
	app.Post("/vibe/match", handler.Match)
	return app, handler
}

func TestVibeHandler_Status(t *testing.T) {
	app, _ := setupVibeApp()

	req := httptest.NewRequest("GET", "/vibe/status", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestVibeHandler_CurrentTrack(t *testing.T) {
	app, _ := setupVibeApp()

	req := httptest.NewRequest("GET", "/vibe/current-track", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestVibeHandler_Match(t *testing.T) {
	app, _ := setupVibeApp()

	body := `{"track_id": "track1"}`
	req := httptest.NewRequest("POST", "/vibe/match", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
}

func TestVibeHandler_Match_InvalidBody(t *testing.T) {
	app, _ := setupVibeApp()

	body := `invalid json`
	req := httptest.NewRequest("POST", "/vibe/match", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}
