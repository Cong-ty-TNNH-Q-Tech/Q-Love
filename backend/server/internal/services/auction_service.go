// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type AuctionService interface {
	StartDailyAuctions(ctx context.Context) error
	PlaceBid(ctx context.Context, auctionID, bidderID uuid.UUID, amount float64) error
	FinalizeAuctions(ctx context.Context) error
	GetActiveAuctions(ctx context.Context) ([]models.BlindAuction, error)
}

type auctionService struct {
	auctionRepo   repository.AuctionRepository
	walletRepo    repository.WalletRepository
	txManager     repository.TransactionManager
	chatLockRepo  repository.ChatLockRepository
}

func NewAuctionService(
	auctionRepo repository.AuctionRepository,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
	chatLockRepo repository.ChatLockRepository,
) AuctionService {
	return &auctionService{
		auctionRepo:   auctionRepo,
		walletRepo:    walletRepo,
		txManager:     txManager,
		chatLockRepo:  chatLockRepo,
	}
}

func (s *auctionService) GetActiveAuctions(ctx context.Context) ([]models.BlindAuction, error) {
	return s.auctionRepo.GetActiveAuctions(ctx)
}

// StartDailyAuctions picks Top 5 users based on some metric (e.g. cards) and creates Blind Auctions.
func (s *auctionService) StartDailyAuctions(ctx context.Context) error {
	// Mock: Fetch top 5 users. In reality, we'd query users ordered by score.
	var topUsers []uuid.UUID
	for i := 0; i < 5; i++ {
		topUsers = append(topUsers, uuid.New())
	}

	for _, uid := range topUsers {
		auction := &models.BlindAuction{
			ID:           uuid.New(),
			TargetUserID: uid,
			StartTime:    time.Now(),
			EndTime:      time.Now().Add(24 * time.Hour),
			Status:       "active",
		}
		if err := s.auctionRepo.CreateAuction(ctx, auction); err != nil {
			return err
		}
	}
	return nil
}

func (s *auctionService) PlaceBid(ctx context.Context, auctionID, bidderID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return errors.New("bid amount must be greater than zero")
	}

	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Check auction
		auction, err := s.auctionRepo.GetAuctionForUpdate(txCtx, auctionID)
		if err != nil {
			return err
		}
		if auction.Status != "active" || time.Now().After(auction.EndTime) {
			return errors.New("auction is not active")
		}
		if auction.TargetUserID == bidderID {
			return errors.New("cannot bid on yourself")
		}

		// 2. Check Wallet
		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, bidderID)
		if err != nil {
			return err
		}

		// Calculate total already bid in this auction by this user
		bids, err := s.auctionRepo.GetBidsByAuction(txCtx, auctionID)
		if err != nil {
			return err
		}
		var maxBid float64
		for _, b := range bids {
			if b.BidderID == bidderID && b.Amount > maxBid {
				maxBid = b.Amount
			}
		}

		// Amount is the additional bid
		if wallet.Balance < amount {
			return errors.New("insufficient balance for this bid")
		}

		// 3. Deduct balance (escrow)
		if err := s.walletRepo.UpdateBalance(txCtx, bidderID, -amount); err != nil {
			return err
		}

		// 4. Create transaction
		walletTx := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      bidderID,
			Amount:      -amount,
			Type:        "AUCTION_BID",
			ReferenceID: auctionID,
		}
		if err := s.walletRepo.CreateTransaction(txCtx, walletTx); err != nil {
			return err
		}

		// 5. Save Bid
		bid := &models.AuctionBid{
			ID:        uuid.New(),
			AuctionID: auctionID,
			BidderID:  bidderID,
			Amount:    maxBid + amount, // Record max bid so far for simplicity
		}
		return s.auctionRepo.PlaceBid(txCtx, bid)
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (s *auctionService) FinalizeAuctions(ctx context.Context) error {
	auctions, err := s.auctionRepo.GetActiveAuctions(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, auction := range auctions {
		if now.Before(auction.EndTime) {
			continue // Not yet ended
		}

		err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			// Lock auction
			lockedAuc, err := s.auctionRepo.GetAuctionForUpdate(txCtx, auction.ID)
			if err != nil || lockedAuc.Status != "active" {
				return err
			}

			bids, err := s.auctionRepo.GetBidsByAuction(txCtx, auction.ID)
			if err != nil {
				return err
			}

			highest, err := s.auctionRepo.GetHighestBid(txCtx, auction.ID)
			if err != nil {
				return err
			}

			if highest == nil {
				// No bids
				return s.auctionRepo.UpdateAuctionStatus(txCtx, auction.ID, "completed", nil, 0)
			}

			// Refund losers
			// To simplify, we keep track of max bid per user
			userMaxBid := make(map[uuid.UUID]float64)
			for _, b := range bids {
				if b.Amount > userMaxBid[b.BidderID] {
					userMaxBid[b.BidderID] = b.Amount
				}
			}

			for uid, amt := range userMaxBid {
				if uid != highest.BidderID {
					// Refund
					if err := s.walletRepo.UpdateBalance(txCtx, uid, amt); err != nil {
						return err
					}
					// Refund transaction
					refundTx := &models.WalletTransaction{
						ID:          uuid.New(),
						UserID:      uid,
						Amount:      amt,
						Type:        "AUCTION_REFUND",
						ReferenceID: auction.ID,
					}
					if err := s.walletRepo.CreateTransaction(txCtx, refundTx); err != nil {
						return err
					}
				}
			}

			// Split 50-50 for winner bid
			half := highest.Amount / 2
			if err := s.walletRepo.UpdateBalance(txCtx, auction.TargetUserID, half); err != nil {
				return err
			}
			commissionTx := &models.WalletTransaction{
				ID:          uuid.New(),
				UserID:      auction.TargetUserID,
				Amount:      half,
				Type:        "AUCTION_REVENUE",
				ReferenceID: auction.ID,
			}
			if err := s.walletRepo.CreateTransaction(txCtx, commissionTx); err != nil {
				return err
			}

			// Create ChatLock
			if s.chatLockRepo != nil {
				lock := models.ChatLock{
					ID:        uuid.New(),
					UserID1:   auction.TargetUserID,
					UserID2:   highest.BidderID,
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}
				if err := s.chatLockRepo.Create(txCtx, &lock); err != nil {
					return err
				}
			}

			// Update Auction Status
			return s.auctionRepo.UpdateAuctionStatus(txCtx, auction.ID, "completed", &highest.BidderID, highest.Amount)
		}, &sql.TxOptions{Isolation: sql.LevelSerializable})

		if err != nil {
			// Log error and continue with next auction
			continue
		}
	}
	return nil
}
