// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardSteal struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	AttackerID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"attacker_id"`
	DefenderID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"defender_id"`
	TargetCardID uuid.UUID      `gorm:"type:uuid;not null" json:"target_card_id"`
	Result       string         `gorm:"type:varchar(20);default:'pending'" json:"result"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
