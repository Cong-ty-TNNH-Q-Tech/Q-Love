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
	"gorm.io/gorm"
)

type AuctionService interface {
	StartDailyAuctions(ctx context.Context) error
	PlaceBid(ctx context.Context, auctionID, bidderID uuid.UUID, amount float64) error
	FinalizeAuctions(ctx context.Context) error
}

type auctionService struct {
	auctionRepo repository.AuctionRepository
	walletRepo  repository.WalletRepository
	txManager   repository.TransactionManager
	userRepo    repository.UserRepository
	db          *gorm.DB // Using db for chat lock and top users for simplicity since repo methods are missing
}

func NewAuctionService(
	auctionRepo repository.AuctionRepository,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
	userRepo repository.UserRepository,
	db *gorm.DB,
) AuctionService {
	return &auctionService{
		auctionRepo: auctionRepo,
		walletRepo:  walletRepo,
		txManager:   txManager,
		userRepo:    userRepo,
		db:          db,
	}
}

func (s *auctionService) StartDailyAuctions(ctx context.Context) error {
	// Query real top users
	topUsers, err := s.userRepo.GetTopUsersByScore(ctx, 5)
	if err != nil {
		return err
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
	limit := 100
	offset := 0
	now := time.Now()

	for {
		auctions, err := s.auctionRepo.GetActiveAuctions(ctx, offset, limit)
		if err != nil {
			return err
		}
		if len(auctions) == 0 {
			break
		}

		var endedAuctionIDs []uuid.UUID
		var endedAuctions []models.BlindAuction
		for _, auction := range auctions {
			if now.After(auction.EndTime) || now.Equal(auction.EndTime) {
				endedAuctionIDs = append(endedAuctionIDs, auction.ID)
				endedAuctions = append(endedAuctions, auction)
			}
		}

		var successfullyUpdated int
		if len(endedAuctionIDs) > 0 {
			bids, err := s.auctionRepo.GetBidsForAuctions(ctx, endedAuctionIDs)
			if err != nil {
				return err
			}

			bidsByAuction := make(map[uuid.UUID][]models.AuctionBid)
			for _, bid := range bids {
				bidsByAuction[bid.AuctionID] = append(bidsByAuction[bid.AuctionID], bid)
			}

			for _, auction := range endedAuctions {
				auctionBids := bidsByAuction[auction.ID]
				
				var highestBid *models.AuctionBid
				for i, b := range auctionBids {
					if highestBid == nil || b.Amount > highestBid.Amount {
						highestBid = &auctionBids[i]
					}
				}

				err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
					lockedAuc, err := s.auctionRepo.GetAuctionForUpdate(txCtx, auction.ID)
					if err != nil || lockedAuc.Status != "active" {
						return err
					}

					if highestBid == nil {
						return s.auctionRepo.UpdateAuctionStatus(txCtx, auction.ID, "completed", nil, 0)
					}

					userMaxBid := make(map[uuid.UUID]float64)
					for _, b := range auctionBids {
						if b.Amount > userMaxBid[b.BidderID] {
							userMaxBid[b.BidderID] = b.Amount
						}
					}

					for uid, amt := range userMaxBid {
						if uid != highestBid.BidderID {
							if err := s.walletRepo.UpdateBalance(txCtx, uid, amt); err != nil {
								return err
							}
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

					half := highestBid.Amount / 2
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

					if s.db != nil {
						lock := models.ChatLock{
							ID:        uuid.New(),
							UserID1:   auction.TargetUserID,
							UserID2:   highestBid.BidderID,
							ExpiresAt: time.Now().Add(24 * time.Hour),
						}
						if err := repository.GetDB(txCtx, s.db).WithContext(txCtx).Create(&lock).Error; err != nil {
							return err
						}
					}

					return s.auctionRepo.UpdateAuctionStatus(txCtx, auction.ID, "completed", &highestBid.BidderID, highestBid.Amount)
				}, &sql.TxOptions{Isolation: sql.LevelSerializable})
				
				if err == nil {
					successfullyUpdated++
				}
			}
		}

		// Calculate new offset:
		// Any auction that was NOT successfully updated to 'completed' (either because it hasn't ended yet, or it failed)
		// will remain in the 'active' status. Therefore, it will appear at the beginning of the next query.
		// To skip over these remaining active auctions, we increment the offset by the number of such auctions.
		remainingActive := len(auctions) - successfullyUpdated
		offset += remainingActive
	}
	return nil
}
