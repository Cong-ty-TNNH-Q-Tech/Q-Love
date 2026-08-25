// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

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
