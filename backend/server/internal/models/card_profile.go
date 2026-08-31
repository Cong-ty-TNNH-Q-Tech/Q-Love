// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type CardProfile struct {
	UserID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	CurrentPrice       float64   `gorm:"type:numeric(15,2);default:100" json:"current_price"`
	TotalCards         int       `gorm:"default:1000" json:"total_cards"`
	AvailableCards     int       `gorm:"default:1000" json:"available_cards"`
	MatchCountCached   int       `gorm:"default:0" json:"-"`
	LocketCountCached  int       `gorm:"default:0" json:"-"`
	ClanUpvoteCached   int       `gorm:"default:0" json:"-"`
	CourtPenaltyCached int       `gorm:"default:0" json:"-"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Virtual fields for API response
	Ticker   string  `gorm:"-" json:"ticker"`
	IsHalted *bool   `gorm:"-" json:"is_halted,omitempty"`
}
