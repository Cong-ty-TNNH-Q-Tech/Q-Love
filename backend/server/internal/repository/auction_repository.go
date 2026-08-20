// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type AuctionRepository interface {
	CreateAuction(ctx context.Context, auction *models.BlindAuction) error
	GetActiveAuctions(ctx context.Context, offset, limit int) ([]models.BlindAuction, error)
	GetActiveAuctionsCursor(ctx context.Context, lastID uuid.UUID, limit int) ([]models.BlindAuction, error)
	GetAuctionForUpdate(ctx context.Context, auctionID uuid.UUID) (*models.BlindAuction, error)
	PlaceBid(ctx context.Context, bid *models.AuctionBid) error
	GetHighestBid(ctx context.Context, auctionID uuid.UUID) (*models.AuctionBid, error)
	GetBidsByAuction(ctx context.Context, auctionID uuid.UUID) ([]models.AuctionBid, error)
	GetBidsForAuctions(ctx context.Context, auctionIDs []uuid.UUID) ([]models.AuctionBid, error)
	UpdateAuctionStatus(ctx context.Context, auctionID uuid.UUID, status string, winnerID *uuid.UUID, winningBid float64) error
}

type auctionRepository struct {
	db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) AuctionRepository {
	return &auctionRepository{db: db}
}

func (r *auctionRepository) CreateAuction(ctx context.Context, auction *models.BlindAuction) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(auction).Error
}

func (r *auctionRepository) GetActiveAuctions(ctx context.Context, offset, limit int) ([]models.BlindAuction, error) {
	var auctions []models.BlindAuction
	err := GetDB(ctx, r.db).WithContext(ctx).
		Where("status = ?", "active").
		Offset(offset).Limit(limit).
		Find(&auctions).Error
	return auctions, err
}

func (r *auctionRepository) GetActiveAuctionsCursor(ctx context.Context, lastID uuid.UUID, limit int) ([]models.BlindAuction, error) {
	var auctions []models.BlindAuction
	query := GetDB(ctx, r.db).WithContext(ctx).Where("status = ?", "active").Order("id ASC").Limit(limit)
	if lastID != uuid.Nil {
		query = query.Where("id > ?", lastID)
	}
	err := query.Find(&auctions).Error
	return auctions, err
}

func (r *auctionRepository) GetAuctionForUpdate(ctx context.Context, auctionID uuid.UUID) (*models.BlindAuction, error) {
	var auction models.BlindAuction
	err := GetDB(ctx, r.db).WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", auctionID).
		First(&auction).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func (r *auctionRepository) PlaceBid(ctx context.Context, bid *models.AuctionBid) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(bid).Error
}

func (r *auctionRepository) GetHighestBid(ctx context.Context, auctionID uuid.UUID) (*models.AuctionBid, error) {
	var bid models.AuctionBid
	err := GetDB(ctx, r.db).WithContext(ctx).
		Where("auction_id = ?", auctionID).
		Order("amount DESC, created_at ASC").
		First(&bid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // no bids
		}
		return nil, err
	}
	return &bid, nil
}

func (r *auctionRepository) GetBidsByAuction(ctx context.Context, auctionID uuid.UUID) ([]models.AuctionBid, error) {
	var bids []models.AuctionBid
	err := GetDB(ctx, r.db).WithContext(ctx).Where("auction_id = ?", auctionID).Find(&bids).Error
	return bids, err
}

func (r *auctionRepository) GetBidsForAuctions(ctx context.Context, auctionIDs []uuid.UUID) ([]models.AuctionBid, error) {
	var bids []models.AuctionBid
	if len(auctionIDs) == 0 {
		return bids, nil
	}
	err := GetDB(ctx, r.db).WithContext(ctx).Where("auction_id IN ?", auctionIDs).Find(&bids).Error
	return bids, err
}

func (r *auctionRepository) UpdateAuctionStatus(ctx context.Context, auctionID uuid.UUID, status string, winnerID *uuid.UUID, winningBid float64) error {
	return GetDB(ctx, r.db).WithContext(ctx).
		Model(&models.BlindAuction{}).
		Where("id = ?", auctionID).
		Updates(map[string]interface{}{
			"status":      status,
			"winner_id":   winnerID,
			"winning_bid": winningBid,
		}).Error
}
