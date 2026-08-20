// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockFeedService struct {
	feed []models.FeedUserResponse
	err  error
}

func (m *mockFeedService) GetFeed(ctx context.Context, userID uuid.UUID, filter string, radius int) ([]models.FeedUserResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.feed, nil
}

func TestFeedHandler_GetFeed_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockFeedService{
		feed: []models.FeedUserResponse{
			{User: models.User{ID: uuid.New(), DOB: time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC)}},
		},
	}
	h := NewFeedHandler(mockSvc)

	app.Get("/feed", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return h.GetFeed(c)
	})

	req := httptest.NewRequest("GET", "/feed?filter=spiritual&radius=10000", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestFeedHandler_GetFeed_Unauthorized(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockFeedService{}
	h := NewFeedHandler(mockSvc)

	app.Get("/feed", func(c *fiber.Ctx) error {
		// user_id is NOT set in Locals
		return h.GetFeed(c)
	})

	req := httptest.NewRequest("GET", "/feed", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestFeedHandler_GetFeed_ServiceError(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockFeedService{
		err: errors.New("db error"),
	}
	h := NewFeedHandler(mockSvc)

	app.Get("/feed", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New())
		return h.GetFeed(c)
	})

	req := httptest.NewRequest("GET", "/feed", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}
