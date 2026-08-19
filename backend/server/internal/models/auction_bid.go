// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package models

import (
	"time"

	"github.com/google/uuid"
)

// AuctionBid represents a bid placed by a user in a blind auction.
type AuctionBid struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AuctionID uuid.UUID `gorm:"type:uuid;not null" json:"auction_id"`
	BidderID  uuid.UUID `gorm:"type:uuid;not null" json:"bidder_id"`
	Amount    float64   `gorm:"not null" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
