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

type mockClanService struct {
	err  error
	clan *models.Clan
}

func (m *mockClanService) CreateClan(ctx context.Context, userID uuid.UUID, name, slogan, logoURL string) (*models.Clan, error) {
	return m.clan, m.err
}

func TestCreateClanHandler(t *testing.T) {
	app := fiber.New()
	
	t.Run("Success", func(t *testing.T) {
		mSvc := &mockClanService{
			clan: &models.Clan{
				Name: "New Clan",
			},
		}
		h := NewClanHandler(mSvc)
		
		appReq := fiber.New()
		appReq.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", uuid.New())
			return c.Next()
		})
		appReq.Post("/clans", h.CreateClan)

		payload := CreateClanRequest{
			Name: "New Clan",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/clans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := appReq.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Missing Name", func(t *testing.T) {
		h := NewClanHandler(&mockClanService{})
		app.Post("/clans_missing", h.CreateClan)

		payload := CreateClanRequest{
			Slogan: "No name",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/clans_missing", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Name Taken", func(t *testing.T) {
		mSvc := &mockClanService{err: errors.New("ERR_CLAN_NAME_TAKEN")}
		h := NewClanHandler(mSvc)
		
		appReq := fiber.New()
		appReq.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", uuid.New())
			return c.Next()
		})
		appReq.Post("/clans", h.CreateClan)

		payload := CreateClanRequest{Name: "Taken"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/clans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := appReq.Test(req, -1)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Insufficient Balance", func(t *testing.T) {
		mSvc := &mockClanService{err: errors.New("insufficient balance")}
		h := NewClanHandler(mSvc)
		
		appReq := fiber.New()
		appReq.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", uuid.New())
			return c.Next()
		})
		appReq.Post("/clans", h.CreateClan)

		payload := CreateClanRequest{Name: "Broke"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/clans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := appReq.Test(req, -1)
		assert.Equal(t, 403, resp.StatusCode)
	})
	
	t.Run("No Auth Context", func(t *testing.T) {
		h := NewClanHandler(&mockClanService{})
		app.Post("/clans_no_auth", h.CreateClan)

		payload := CreateClanRequest{Name: "NoAuth"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/clans_no_auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		assert.Equal(t, 401, resp.StatusCode)
	})
}
