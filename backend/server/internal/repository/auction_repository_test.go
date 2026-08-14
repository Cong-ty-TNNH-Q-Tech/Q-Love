// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

func TestAuctionRepository_CreateAuction(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	auction := &models.BlindAuction{
		ID:           uuid.New(),
		TargetUserID: uuid.New(),
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(24 * time.Hour),
		Status:       "active",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "blind_auctions"`)).
		WithArgs(auction.ID, auction.TargetUserID, auction.StartTime, auction.EndTime, auction.Status, sqlmock.AnyArg(), auction.WinningBid, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateAuction(context.Background(), auction)
	assert.NoError(t, err)
}

func TestAuctionRepository_GetActiveAuctions(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "blind_auctions" WHERE status = $1`)).
		WithArgs("active").
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(uuid.New(), "active"))

	auctions, err := repo.GetActiveAuctions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, auctions, 1)
}

func TestAuctionRepository_PlaceBid(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	bid := &models.AuctionBid{
		ID:        uuid.New(),
		AuctionID: uuid.New(),
		BidderID:  uuid.New(),
		Amount:    1500,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auction_bids"`)).
		WithArgs(bid.ID, bid.AuctionID, bid.BidderID, bid.Amount, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.PlaceBid(context.Background(), bid)
	assert.NoError(t, err)
}

func TestAuctionRepository_UpdateAuctionStatus(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	auctionID := uuid.New()
	winnerID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "blind_auctions" SET`)).
		WithArgs("completed", winnerID, float64(2000), sqlmock.AnyArg(), auctionID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateAuctionStatus(context.Background(), auctionID, "completed", &winnerID, 2000)
	assert.NoError(t, err)
}

func TestAuctionRepository_GetAuctionForUpdate(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	auctionID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "blind_auctions" WHERE id = $1 FOR UPDATE`)).
		WithArgs(auctionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(auctionID, "active"))

	auction, err := repo.GetAuctionForUpdate(context.Background(), auctionID)
	assert.NoError(t, err)
	assert.NotNil(t, auction)
	assert.Equal(t, auctionID, auction.ID)
}

func TestAuctionRepository_GetHighestBid(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	auctionID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auction_bids" WHERE auction_id = $1 ORDER BY amount desc LIMIT $2`)).
		WithArgs(auctionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "amount"}).AddRow(uuid.New(), auctionID, float64(2000)))

	bid, err := repo.GetHighestBid(context.Background(), auctionID)
	assert.NoError(t, err)
	assert.NotNil(t, bid)
	assert.Equal(t, float64(2000), bid.Amount)
}

func TestAuctionRepository_GetBidsByAuction(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := NewAuctionRepository(db)
	auctionID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auction_bids" WHERE auction_id = $1`)).
		WithArgs(auctionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "amount"}).AddRow(uuid.New(), auctionID, float64(1500)).AddRow(uuid.New(), auctionID, float64(2000)))

	bids, err := repo.GetBidsByAuction(context.Background(), auctionID)
	assert.NoError(t, err)
	assert.Len(t, bids, 2)
}
