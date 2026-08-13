// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type mockAuctionRepo struct {
	auction *models.BlindAuction
	bids    []models.AuctionBid
	highest *models.AuctionBid
	err     error
}

func (m *mockAuctionRepo) CreateAuction(ctx context.Context, auction *models.BlindAuction) error {
	return m.err
}
func (m *mockAuctionRepo) GetActiveAuctions(ctx context.Context) ([]models.BlindAuction, error) {
	if m.auction != nil {
		return []models.BlindAuction{*m.auction}, m.err
	}
	return nil, m.err
}
func (m *mockAuctionRepo) GetAuctionForUpdate(ctx context.Context, auctionID uuid.UUID) (*models.BlindAuction, error) {
	return m.auction, m.err
}
func (m *mockAuctionRepo) PlaceBid(ctx context.Context, bid *models.AuctionBid) error {
	return m.err
}
func (m *mockAuctionRepo) GetHighestBid(ctx context.Context, auctionID uuid.UUID) (*models.AuctionBid, error) {
	return m.highest, m.err
}
func (m *mockAuctionRepo) GetBidsByAuction(ctx context.Context, auctionID uuid.UUID) ([]models.AuctionBid, error) {
	return m.bids, m.err
}
func (m *mockAuctionRepo) UpdateAuctionStatus(ctx context.Context, auctionID uuid.UUID, status string, winnerID *uuid.UUID, winningBid float64) error {
	return m.err
}

func TestAuctionService_PlaceBid_Success(t *testing.T) {
	auctionID := uuid.New()
	bidderID := uuid.New()
	targetID := uuid.New()

	auction := &models.BlindAuction{
		ID:           auctionID,
		TargetUserID: targetID,
		StartTime:    time.Now().Add(-1 * time.Hour),
		EndTime:      time.Now().Add(23 * time.Hour),
		Status:       "active",
	}

	walletRepo := &mockMinigameWalletRepo{
		wallet: &models.UserWallet{UserID: bidderID, Balance: 1000},
	}
	auctionRepo := &mockAuctionRepo{auction: auction}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, nil)
	err := service.PlaceBid(context.Background(), auctionID, bidderID, 500)
	assert.NoError(t, err)
}

func TestAuctionService_FinalizeAuctions(t *testing.T) {
	auctionID := uuid.New()
	bidderID := uuid.New()
	targetID := uuid.New()

	auction := &models.BlindAuction{
		ID:           auctionID,
		TargetUserID: targetID,
		StartTime:    time.Now().Add(-25 * time.Hour),
		EndTime:      time.Now().Add(-1 * time.Hour), // Ended
		Status:       "active",
	}

	bid := models.AuctionBid{
		ID:        uuid.New(),
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    1000,
	}

	walletRepo := &mockMinigameWalletRepo{
		wallet: &models.UserWallet{UserID: bidderID, Balance: 0},
	}
	auctionRepo := &mockAuctionRepo{
		auction: auction,
		bids:    []models.AuctionBid{bid},
		highest: &bid,
	}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, nil)
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err)
}
