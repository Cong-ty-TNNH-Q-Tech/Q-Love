// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type Clan struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name               string    `json:"name" gorm:"type:varchar(100);unique;not null"`
	Slogan             string    `json:"slogan" gorm:"type:varchar(255)"`
	LogoURL            string    `json:"logo_url" gorm:"type:text"`
	LeaderID           uuid.UUID `json:"leader_id" gorm:"type:uuid;not null"`
	WeeklyScore        int       `json:"weekly_score" gorm:"default:0"`
	CampfireStreak     int       `json:"campfire_streak" gorm:"default:0"`
	DailyActiveMembers int       `json:"daily_active_members" gorm:"default:0"`
	LastCampfireAt     time.Time `json:"last_campfire_at"`
	CreatedAt          time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"default:now()"`
}

type ClanMember struct {
	ClanID   uuid.UUID `json:"clan_id" gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	Role     string    `json:"role" gorm:"type:varchar(20);default:'member'"` // 'leader', 'member'
	JoinedAt time.Time `json:"joined_at" gorm:"default:now()"`
}
