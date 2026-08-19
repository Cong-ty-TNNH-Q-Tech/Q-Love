// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VibeMatch struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserA     uuid.UUID      `gorm:"type:uuid;not null" json:"user_a"`
	UserB     uuid.UUID      `gorm:"type:uuid;not null" json:"user_b"`
	TrackID   string         `gorm:"type:varchar(255);not null" json:"track_id"`
	RoomID    string         `gorm:"type:varchar(255);not null" json:"room_id"`
	Status    string         `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
