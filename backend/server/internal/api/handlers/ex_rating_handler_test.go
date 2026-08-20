// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockExRatingService struct{}

func (m *mockExRatingService) SubmitRating(ctx context.Context, targetUserID, matchID uuid.UUID, score int, tags []string) error {
	return nil
}
func (m *mockExRatingService) ViewRating(ctx context.Context, viewerID, targetUserID uuid.UUID) (float64, int64, map[string]int, error) {
	return 4.5, 10, map[string]int{"#green": 5}, nil
}

func TestExRatingHandler_SubmitRating(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingService{})
	app.Post("/ex-ratings", handler.SubmitRating)

	reqBody := map[string]interface{}{
		"target_user_id": uuid.New().String(),
		"match_id":       uuid.New().String(),
		"rating_score":   5,
		"tags":           []string{"#tốt"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestExRatingHandler_ViewRating(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingService{})
	app.Get("/users/:user_id/ex-rating", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.ViewRating(c)
	})

	targetID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/users/"+targetID+"/ex-rating", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, float64(4.5), result["average_rating"])
func TestExRatingHandler_SubmitRating_InvalidBody(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingService{})
	app.Post("/ex-ratings", handler.SubmitRating)

	req := httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestExRatingHandler_SubmitRating_InvalidScore(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingService{})
	app.Post("/ex-ratings", handler.SubmitRating)

	jsonBody, _ := json.Marshal(map[string]interface{}{"rating_score": 6})
	req := httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestExRatingHandler_SubmitRating_InvalidIDs(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingService{})
	app.Post("/ex-ratings", handler.SubmitRating)

	// Invalid target user ID
	jsonBody, _ := json.Marshal(map[string]interface{}{"rating_score": 5, "target_user_id": "invalid"})
	req := httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Invalid match ID
	jsonBody, _ = json.Marshal(map[string]interface{}{"rating_score": 5, "target_user_id": uuid.New().String(), "match_id": "invalid"})
	req = httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

type mockExRatingServiceErr struct{}
func (m *mockExRatingServiceErr) SubmitRating(ctx context.Context, targetUserID, matchID uuid.UUID, score int, tags []string) error {
	return fiber.ErrBadRequest
}
func (m *mockExRatingServiceErr) ViewRating(ctx context.Context, viewerID, targetUserID uuid.UUID) (float64, int64, map[string]int, error) {
	return 0, 0, nil, fiber.ErrBadRequest
}

func TestExRatingHandler_SubmitRating_ServiceErr(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingServiceErr{})
	app.Post("/ex-ratings", handler.SubmitRating)

	jsonBody, _ := json.Marshal(map[string]interface{}{"rating_score": 5, "target_user_id": uuid.New().String(), "match_id": uuid.New().String()})
	req := httptest.NewRequest(http.MethodPost, "/ex-ratings", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestExRatingHandler_ViewRating_Errors(t *testing.T) {
	app := fiber.New()
	handler := NewExRatingHandler(&mockExRatingServiceErr{})
	app.Get("/users/:user_id/ex-rating", func(c *fiber.Ctx) error {
		val := c.Query("locals")
		if val == "none" {
			// do nothing
		} else if val == "invalid" {
			c.Locals("user_id", "invalid-uuid")
		} else {
			c.Locals("user_id", uuid.New().String())
		}
		return handler.ViewRating(c)
	})

	// 1. Missing locals
	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/ex-rating?locals=none", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	// 2. Invalid locals UUID
	req = httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/ex-rating?locals=invalid", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	// 3. Invalid target ID
	req = httptest.NewRequest(http.MethodGet, "/users/invalid/ex-rating", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// 4. Service Error
	req = httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String()+"/ex-rating", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
