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

type mockCourtService struct {
	err error
}

func (m *mockCourtService) FileLawsuit(ctx context.Context, plaintiffID, defendantID, matchID uuid.UUID, reason string) (*models.CourtCase, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &models.CourtCase{}, nil
}

func (m *mockCourtService) GetFeed(ctx context.Context, jurorID uuid.UUID, limit int) ([]models.CourtCase, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []models.CourtCase{}, nil
}

func (m *mockCourtService) VoteCase(ctx context.Context, caseID, jurorID uuid.UUID, voteType models.CourtVoteType) error {
	return m.err
}

func (m *mockCourtService) WithdrawCase(ctx context.Context, caseID, plaintiffID uuid.UUID) error {
	return m.err
}

func setupCourtHandlerApp(mockSvc *mockCourtService) *fiber.App {
	app := fiber.New()
	h := NewCourtHandler(mockSvc)
	app.Post("/court/cases", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		return h.FileLawsuit(c)
	})
	app.Post("/court/:case_id/vote", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		return h.VoteCase(c)
	})
	app.Post("/court/:case_id/withdraw", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		return h.WithdrawCase(c)
	})
	app.Get("/court/feed", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New().String())
		return h.GetFeed(c)
	})
	return app
}

func TestCourtHandler_FileCase_Success(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	
	body := map[string]interface{}{
		"defendant_id": uuid.New().String(),
		"match_id": uuid.New().String(),
		"reason": "Test reason",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/court/cases", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestCourtHandler_FileCase_BadRequest(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	req := httptest.NewRequest("POST", "/court/cases", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCourtHandler_FileCase_ServiceErr(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{err: errors.New("cannot file case")})
	body := map[string]interface{}{
		"defendant_id": uuid.New().String(),
		"match_id": uuid.New().String(),
		"reason": "Test reason",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/court/cases", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCourtHandler_VoteCase_Success(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	body := map[string]interface{}{
		"vote": "guilty",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/court/"+uuid.New().String()+"/vote", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCourtHandler_WithdrawCase_Success(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	req := httptest.NewRequest("POST", "/court/"+uuid.New().String()+"/withdraw", nil)
	
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCourtHandler_GetFeed_Success(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	req := httptest.NewRequest("GET", "/court/feed", nil)
	
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCourtHandler_GetFeed_ServiceErr(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{err: errors.New("error")})
	req := httptest.NewRequest("GET", "/court/feed", nil)
	
	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestCourtHandler_VoteCase_BadRequest(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	req := httptest.NewRequest("POST", "/court/"+uuid.New().String()+"/vote", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCourtHandler_VoteCase_InvalidVote(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{})
	body := map[string]interface{}{
		"vote": "invalid_vote",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/court/"+uuid.New().String()+"/vote", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCourtHandler_WithdrawCase_ServiceErr(t *testing.T) {
	app := setupCourtHandlerApp(&mockCourtService{err: errors.New("cannot withdraw")})
	req := httptest.NewRequest("POST", "/court/"+uuid.New().String()+"/withdraw", nil)
	
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}
