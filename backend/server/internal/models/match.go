// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	User1ID            uuid.UUID      `gorm:"type:uuid;not null" json:"user1_id"`
	User2ID            uuid.UUID      `gorm:"type:uuid;not null" json:"user2_id"`
	StreakScore        int            `gorm:"default:0" json:"streak_score"`
	HighestStreakScore int            `gorm:"default:0" json:"highest_streak_score"`
	LastInteractionAt  time.Time      `json:"last_interaction_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
