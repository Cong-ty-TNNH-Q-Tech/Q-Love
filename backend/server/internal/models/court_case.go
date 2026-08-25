// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourtCase struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PlaintiffID uuid.UUID      `gorm:"type:uuid;not null;index" json:"plaintiff_id"`
	DefendantID uuid.UUID      `gorm:"type:uuid;not null;index" json:"defendant_id"`
	Status      string         `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
