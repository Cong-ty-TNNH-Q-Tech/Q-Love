// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	MatchID   uuid.UUID `gorm:"type:uuid;not null"`
	SenderID  uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    `gorm:"type:varchar;not null"`
	Content   string    `gorm:"type:text"`
	BlurURL   string    `gorm:"type:text"`
	BlurLevel int       `gorm:"type:int"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
