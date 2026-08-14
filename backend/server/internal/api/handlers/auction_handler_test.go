// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

// mockAuctionService and mockAuctionRepository
type mockAuctionService struct {
	err error
}
func (m *mockAuctionService) StartDailyAuctions(ctx context.Context) error { return m.err }
func (m *mockAuctionService) PlaceBid(ctx context.Context, auctionID, bidderID uuid.UUID, amount float64) error { return m.err }
func (m *mockAuctionService) FinalizeAuctions(ctx context.Context) error { return m.err }

type mockAuctionRepo struct {
	auctions []models.BlindAuction
	err      error
}
func (m *mockAuctionRepo) CreateAuction(ctx context.Context, auction *models.BlindAuction) error { return m.err }
func (m *mockAuctionRepo) GetActiveAuctions(ctx context.Context) ([]models.BlindAuction, error) { return m.auctions, m.err }
func (m *mockAuctionRepo) GetAuctionForUpdate(ctx context.Context, auctionID uuid.UUID) (*models.BlindAuction, error) { return nil, m.err }
func (m *mockAuctionRepo) PlaceBid(ctx context.Context, bid *models.AuctionBid) error { return m.err }
func (m *mockAuctionRepo) GetHighestBid(ctx context.Context, auctionID uuid.UUID) (*models.AuctionBid, error) { return nil, m.err }
func (m *mockAuctionRepo) GetBidsByAuction(ctx context.Context, auctionID uuid.UUID) ([]models.AuctionBid, error) { return nil, m.err }
func (m *mockAuctionRepo) UpdateAuctionStatus(ctx context.Context, auctionID uuid.UUID, status string, winnerID *uuid.UUID, winningBid float64) error { return m.err }


func TestAuctionHandler_GetActiveAuctions(t *testing.T) {
	app := fiber.New()
	svc := &mockAuctionService{}
	repo := &mockAuctionRepo{
		auctions: []models.BlindAuction{
			{ID: uuid.New(), TargetUserID: uuid.New(), Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	h := NewAuctionHandler(svc, repo)
	app.Get("/api/v1/auctions/active", h.GetActiveAuctions)

	req := httptest.NewRequest("GET", "/api/v1/auctions/active", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuctionHandler_PlaceBid(t *testing.T) {
	app := fiber.New()
	svc := &mockAuctionService{}
	repo := &mockAuctionRepo{}
	h := NewAuctionHandler(svc, repo)

	app.Post("/api/v1/auctions/:id/bid", func(c *fiber.Ctx) error {
		// Mock JWT middleware
		c.Locals("userID", uuid.New())
		return h.PlaceBid(c)
	})

	body, _ := json.Marshal(map[string]interface{}{"amount": 1000})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+uuid.New().String()+"/bid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestAuctionHandler_GetActiveAuctions_Error(t *testing.T) {
	app := fiber.New()
	svc := &mockAuctionService{}
	repo := &mockAuctionRepo{
		err: assert.AnError,
	}
	h := NewAuctionHandler(svc, repo)
	app.Get("/api/v1/auctions/active", h.GetActiveAuctions)

	req := httptest.NewRequest("GET", "/api/v1/auctions/active", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestAuctionHandler_PlaceBid_InvalidInput(t *testing.T) {
	app := fiber.New()
	svc := &mockAuctionService{}
	repo := &mockAuctionRepo{}
	h := NewAuctionHandler(svc, repo)

	app.Post("/api/v1/auctions/:id/bid", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New())
		return h.PlaceBid(c)
	})

	// Missing amount
	body, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+uuid.New().String()+"/bid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuctionHandler_PlaceBid_ServiceError(t *testing.T) {
	app := fiber.New()
	svc := &mockAuctionService{err: assert.AnError}
	repo := &mockAuctionRepo{}
	h := NewAuctionHandler(svc, repo)

	app.Post("/api/v1/auctions/:id/bid", func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.New())
		return h.PlaceBid(c)
	})

	body, _ := json.Marshal(map[string]interface{}{"amount": 1000})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+uuid.New().String()+"/bid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}
