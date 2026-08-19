// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatLock represents an exclusive 24-hour chat lock granted to the winner of a blind auction.
type ChatLock struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID1   uuid.UUID `gorm:"type:uuid;not null" json:"user_id_1"` // The auctioned user
	UserID2   uuid.UUID `gorm:"type:uuid;not null" json:"user_id_2"` // The winner
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
