// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount      float64   `gorm:"type:numeric;not null" json:"amount"`
	Type        string    `gorm:"type:varchar(50)" json:"type"` // deposit, contract_hold, penalty, card_trade, wingman_commission
	ReferenceID uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserWallet struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Balance     float64   `gorm:"type:numeric;default:0" json:"balance"`
	HoldBalance float64   `gorm:"type:numeric;default:0" json:"hold_balance"`
}
