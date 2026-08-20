// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDeviceHandler_RegisterFCMToken_Unauthorized(t *testing.T) {
	app := fiber.New()
	handler := NewDeviceHandler(nil)
	app.Post("/devices/token", handler.RegisterFCMToken)

	// Sending request without user_id in context locals (simulating unauthorized)
	body, _ := json.Marshal(map[string]string{"token": "abc"})
	req := httptest.NewRequest("POST", "/devices/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestDeviceHandler_RegisterFCMToken_InvalidPayload(t *testing.T) {
	app := fiber.New()
	handler := NewDeviceHandler(nil)
	app.Post("/devices/token", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return handler.RegisterFCMToken(c)
	})

	// Missing token in body
	body, _ := json.Marshal(map[string]string{"other": "abc"})
	req := httptest.NewRequest("POST", "/devices/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
