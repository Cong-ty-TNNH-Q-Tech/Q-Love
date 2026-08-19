// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package models

import (
	"time"

	"github.com/google/uuid"
)

// BlindAuction represents a daily blind auction for an exclusive 24h chat with a top user.
type BlindAuction struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TargetUserID uuid.UUID `gorm:"type:uuid;not null" json:"target_user_id"`
	StartTime    time.Time `gorm:"not null" json:"start_time"`
	EndTime      time.Time `gorm:"not null" json:"end_time"`
	Status       string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, completed, cancelled
	WinnerID     *uuid.UUID `gorm:"type:uuid" json:"winner_id,omitempty"`
	WinningBid   float64   `gorm:"default:0" json:"winning_bid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
