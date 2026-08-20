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
func (m *mockAuctionRepo) GetActiveAuctions(ctx context.Context, offset, limit int) ([]models.BlindAuction, error) {
	if offset > 0 {
		return nil, m.err
	}
	if m.auction != nil && m.auction.Status == "active" {
		return []models.BlindAuction{*m.auction}, m.err
	}
	return nil, m.err
}
func (m *mockAuctionRepo) GetActiveAuctionsCursor(ctx context.Context, lastID uuid.UUID, limit int) ([]models.BlindAuction, error) {
	if lastID != uuid.Nil {
		return nil, m.err
	}
	if m.auction != nil && m.auction.Status == "active" {
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
func (m *mockAuctionRepo) GetBidsForAuctions(ctx context.Context, auctionIDs []uuid.UUID) ([]models.AuctionBid, error) {
	return m.bids, m.err
}
func (m *mockAuctionRepo) UpdateAuctionStatus(ctx context.Context, auctionID uuid.UUID, status string, winnerID *uuid.UUID, winningBid float64) error {
	if m.auction != nil && m.auction.ID == auctionID {
		m.auction.Status = status
	}
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

type mockUserRepo struct{}
func (m *mockUserRepo) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) {
	var users []uuid.UUID
	for i := 0; i < limit; i++ {
		users = append(users, uuid.New())
	}
	return users, nil
}

type mockChatLockRepo struct {
	err error
}
func (m *mockChatLockRepo) Create(ctx context.Context, lock *models.ChatLock) error {
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

	walletRepo := &mockAuctionWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: bidderID, Balance: 1000}, nil
		},
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error { return nil },
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error { return nil },
	}
	auctionRepo := &mockAuctionRepo{auction: auction}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
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

	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err)
}


func TestAuctionService_StartDailyAuctions(t *testing.T) {
	auctionRepo := &mockAuctionRepo{}
	walletRepo := &mockAuctionWalletRepo{}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.StartDailyAuctions(context.Background())
	// Usually this returns nil if successful or no db is provided in mock
	assert.NoError(t, err)
}

type mockUserRepoError struct{}
func (m *mockUserRepoError) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) {
	return nil, assert.AnError
}

func TestAuctionService_StartDailyAuctions_UserRepoError(t *testing.T) {
	auctionRepo := &mockAuctionRepo{}
	service := NewAuctionService(auctionRepo, nil, nil, &mockUserRepoError{}, nil)
	err := service.StartDailyAuctions(context.Background())
	assert.NoError(t, err)
}

func TestAuctionService_PlaceBid_AmountZeroOrLess(t *testing.T) {
	service := NewAuctionService(nil, nil, nil, &mockUserRepo{}, nil)
	err := service.PlaceBid(context.Background(), uuid.New(), uuid.New(), 0)
	assert.Error(t, err)
	assert.Equal(t, "bid amount must be greater than zero", err.Error())
}

func TestAuctionService_PlaceBid_AuctionNotActive(t *testing.T) {
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{Status: "completed"},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, nil, txManager, &mockUserRepo{}, nil)
	err := service.PlaceBid(context.Background(), uuid.New(), uuid.New(), 100)
	assert.Error(t, err)
	assert.Equal(t, "auction is not active", err.Error())
}

func TestAuctionService_PlaceBid_BidOnSelf(t *testing.T) {
	uid := uuid.New()
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{Status: "active", EndTime: time.Now().Add(1*time.Hour), TargetUserID: uid},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, nil, txManager, &mockUserRepo{}, nil)
	err := service.PlaceBid(context.Background(), uuid.New(), uid, 100)
	assert.Error(t, err)
	assert.Equal(t, "cannot bid on yourself", err.Error())
}

func TestAuctionService_PlaceBid_InsufficientBalance(t *testing.T) {
	uid := uuid.New()
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{Status: "active", EndTime: time.Now().Add(1*time.Hour), TargetUserID: uuid.New()},
	}
	walletRepo := &mockAuctionWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: uid, Balance: 50}, nil
		},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.PlaceBid(context.Background(), uuid.New(), uid, 100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance for this bid", err.Error())
}

func TestAuctionService_PlaceBid_GetWalletError(t *testing.T) {
	uid := uuid.New()
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{Status: "active", EndTime: time.Now().Add(1*time.Hour), TargetUserID: uuid.New()},
	}
	walletRepo := &mockAuctionWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return nil, assert.AnError
		},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.PlaceBid(context.Background(), uuid.New(), uid, 100)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestAuctionService_FinalizeAuctions_GetActiveError(t *testing.T) {
	auctionRepo := &mockAuctionRepo{err: assert.AnError}
	service := NewAuctionService(auctionRepo, nil, nil, &mockUserRepo{}, nil)
	err := service.FinalizeAuctions(context.Background())
	assert.Error(t, err)
}

func TestAuctionService_FinalizeAuctions_AuctionNotEnded(t *testing.T) {
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{EndTime: time.Now().Add(1 * time.Hour)},
	}
	service := NewAuctionService(auctionRepo, nil, nil, &mockUserRepo{}, nil)
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err) 
}

func TestAuctionService_FinalizeAuctions_NoBids(t *testing.T) {
	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{
			ID:      uuid.New(),
			Status:  "active",
			EndTime: time.Now().Add(-1 * time.Hour),
		},
		highest: nil,
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, nil, txManager, &mockUserRepo{}, nil)
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err) 
}

func TestAuctionService_FinalizeAuctions_WithBidsAndRefunds(t *testing.T) {
	winnerID := uuid.New()
	loserID := uuid.New()
	targetID := uuid.New()
	auctionID := uuid.New()

	highestBid := &models.AuctionBid{BidderID: winnerID, Amount: 1000}
	loserBid := models.AuctionBid{BidderID: loserID, Amount: 500}

	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{
			ID:           auctionID,
			TargetUserID: targetID,
			Status:       "active",
			EndTime:      time.Now().Add(-1 * time.Hour),
		},
		bids:    []models.AuctionBid{*highestBid, loserBid},
		highest: highestBid,
	}

	walletRepo := &mockAuctionWalletRepo{}
	txManager := &mockTxManager{}

	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err)
}

func TestAuctionService_StartDailyAuctions_CreateError(t *testing.T) {
	auctionRepo := &mockAuctionRepo{err: assert.AnError}
	service := NewAuctionService(auctionRepo, nil, nil, &mockUserRepo{}, nil)
	err := service.StartDailyAuctions(context.Background())
	assert.Error(t, err)
}

func TestAuctionService_FinalizeAuctions_Refund_UpdateBalance_Error(t *testing.T) {
	winnerID := uuid.New()
	loserID := uuid.New()
	targetID := uuid.New()
	auctionID := uuid.New()

	highestBid := &models.AuctionBid{BidderID: winnerID, Amount: 1000}
	loserBid := models.AuctionBid{BidderID: loserID, Amount: 500}

	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{
			ID:           auctionID,
			TargetUserID: targetID,
			Status:       "active",
			EndTime:      time.Now().Add(-1 * time.Hour),
		},
		bids:    []models.AuctionBid{*highestBid, loserBid},
		highest: highestBid,
	}

	walletRepo := &mockAuctionWalletRepo{
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error {
			if userID == loserID {
				return assert.AnError
			}
			return nil
		},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err) // It continues on inner loop error
}

func TestAuctionService_FinalizeAuctions_Winner_UpdateBalance_Error(t *testing.T) {
	winnerID := uuid.New()
	loserID := uuid.New()
	targetID := uuid.New()
	auctionID := uuid.New()

	highestBid := &models.AuctionBid{BidderID: winnerID, Amount: 1000}
	loserBid := models.AuctionBid{BidderID: loserID, Amount: 500}

	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{
			ID:           auctionID,
			TargetUserID: targetID,
			Status:       "active",
			EndTime:      time.Now().Add(-1 * time.Hour),
		},
		bids:    []models.AuctionBid{*highestBid, loserBid},
		highest: highestBid,
	}

	walletRepo := &mockAuctionWalletRepo{
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error {
			if userID == targetID { // The winner receiving 50-50 split
				return assert.AnError
			}
			return nil
		},
	}
	txManager := &mockTxManager{}
	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{})
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err)
}

func TestAuctionService_FinalizeAuctions_ChatLock_DB_Error(t *testing.T) {
	winnerID := uuid.New()
	loserID := uuid.New()
	targetID := uuid.New()
	auctionID := uuid.New()

	highestBid := &models.AuctionBid{BidderID: winnerID, Amount: 1000}
	loserBid := models.AuctionBid{BidderID: loserID, Amount: 500}

	auctionRepo := &mockAuctionRepo{
		auction: &models.BlindAuction{
			ID:           auctionID,
			TargetUserID: targetID,
			Status:       "active",
			EndTime:      time.Now().Add(-1 * time.Hour),
		},
		bids:    []models.AuctionBid{*highestBid, loserBid},
		highest: highestBid,
	}

	walletRepo := &mockAuctionWalletRepo{}
	txManager := &mockTxManager{}
	
	service := NewAuctionService(auctionRepo, walletRepo, txManager, &mockUserRepo{}, &mockChatLockRepo{err: assert.AnError})
	err := service.FinalizeAuctions(context.Background())
	assert.NoError(t, err)
}

func TestAuctionService_GetActiveAuctions(t *testing.T) {
	auctionID := uuid.New()
	auction := &models.BlindAuction{
		ID:           auctionID,
		Status:       "active",
	}
	auctionRepo := &mockAuctionRepo{auction: auction}
	service := NewAuctionService(auctionRepo, nil, nil, nil, nil)
	
	auctions, err := service.GetActiveAuctions(context.Background(), 0, 100)
	assert.NoError(t, err)
	assert.Len(t, auctions, 1)
	assert.Equal(t, auctionID, auctions[0].ID)
}

func TestAuctionService_GetActiveAuctions_Error(t *testing.T) {
	auctionRepo := &mockAuctionRepo{err: assert.AnError}
	service := NewAuctionService(auctionRepo, nil, nil, nil, nil)
	
	auctions, err := service.GetActiveAuctions(context.Background(), 0, 100)
	assert.Error(t, err)
	assert.Nil(t, auctions)
}

