// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

func TestAuctionRepository(t *testing.T) {
	db := setupTestDB(t)
	err := db.AutoMigrate(&models.BlindAuction{}, &models.AuctionBid{})
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	ctx := context.Background()

	// 1. Create Auction
	auction := &models.BlindAuction{
		ID:           uuid.New(),
		TargetUserID: uuid.New(),
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(24 * time.Hour),
		Status:       "active",
	}
	err = repo.CreateAuction(ctx, auction)
	assert.NoError(t, err)

	// 2. Get Active Auctions
	active, err := repo.GetActiveAuctions(ctx)
	assert.NoError(t, err)
	assert.Len(t, active, 1)

	// 3. Place Bid
	bid := &models.AuctionBid{
		ID:        uuid.New(),
		AuctionID: auction.ID,
		BidderID:  uuid.New(),
		Amount:    1500,
	}
	err = repo.PlaceBid(ctx, bid)
	assert.NoError(t, err)

	bid2 := &models.AuctionBid{
		ID:        uuid.New(),
		AuctionID: auction.ID,
		BidderID:  uuid.New(),
		Amount:    2000,
	}
	err = repo.PlaceBid(ctx, bid2)
	assert.NoError(t, err)

	// 4. Get Highest Bid
	highest, err := repo.GetHighestBid(ctx, auction.ID)
	assert.NoError(t, err)
	assert.NotNil(t, highest)
	assert.Equal(t, float64(2000), highest.Amount)

	// 5. Get Bids by Auction
	bids, err := repo.GetBidsByAuction(ctx, auction.ID)
	assert.NoError(t, err)
	assert.Len(t, bids, 2)

	// 6. Update Status
	err = repo.UpdateAuctionStatus(ctx, auction.ID, "completed", &bid2.BidderID, 2000)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetAuctionForUpdate(ctx, auction.ID)
	assert.NoError(t, err)
	assert.Equal(t, "completed", updated.Status)
	assert.Equal(t, bid2.BidderID, *updated.WinnerID)
	assert.Equal(t, float64(2000), updated.WinningBid)
}
