// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Voucher struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Brand     string         `json:"brand" gorm:"type:varchar(50);not null"`
	Code      string         `json:"code" gorm:"type:varchar(50);unique;not null"`
	ValueXu   int            `json:"value_xu" gorm:"not null"`
	Status    string         `json:"status" gorm:"type:varchar(20);default:'available'"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserVoucher struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	VoucherID uuid.UUID      `json:"voucher_id" gorm:"type:uuid;not null;uniqueIndex"`
	ClaimedAt time.Time      `json:"claimed_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Associations
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Voucher *Voucher `json:"voucher,omitempty" gorm:"foreignKey:VoucherID"`
}
