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

type mockAuctionWalletRepo struct {
	getWalletForUpdateFn func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	updateBalanceFn      func(ctx context.Context, userID uuid.UUID, delta float64) error
	createTransactionFn  func(ctx context.Context, txn *models.WalletTransaction) error
}
func (m *mockAuctionWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockAuctionWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	if m.createTransactionFn != nil { return m.createTransactionFn(ctx, txn) }
	return nil
}
func (m *mockAuctionWalletRepo) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	if m.getWalletForUpdateFn != nil { return m.getWalletForUpdateFn(ctx, userID) }
	return nil, nil
}
func (m *mockAuctionWalletRepo) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	if m.updateBalanceFn != nil { return m.updateBalanceFn(ctx, userID, delta) }
	return nil
}
func (m *mockAuctionWalletRepo) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }

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

	walletRepo := &mockAuctionWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: bidderID, Balance: 1000}, nil
		},
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error { return nil },
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error { return nil },
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

	walletRepo := &mockAuctionWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: bidderID, Balance: 0}, nil
		},
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error { return nil },
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error { return nil },
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


func TestAuctionService_StartDailyAuctions(t *testing.T) {
	auctionRepo := &mockAuctionRepo{}
	walletRepo := &mockAuctionWalletRepo{}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, nil)
	err := service.StartDailyAuctions(context.Background())
	// Usually this returns nil if successful or no db is provided in mock
	assert.NoError(t, err)
}
