// Copyright 2026 Q-Tech Team.
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WallOfShame struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Reason         string         `gorm:"type:text;not null" json:"reason"`
	TomatoesThrown int            `gorm:"type:int;default:0" json:"tomatoes_thrown"`
	ExpiresAt      time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type WallOfShameResponse struct {
	WallOfShame
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`
}
