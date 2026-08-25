// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatingContract struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserAID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_a_id"`
	UserBID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_b_id"`
	DepositAmount   float64        `gorm:"type:numeric;not null" json:"deposit_amount"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, active, completed, cancelled
	CancelledByID   *uuid.UUID     `gorm:"type:uuid" json:"cancelled_by_id,omitempty"`
	TOTPSecret      string         `gorm:"type:varchar(255)" json:"-"`
	AppointmentTime *time.Time     `json:"appointment_time,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
