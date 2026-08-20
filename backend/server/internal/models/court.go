// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourtCaseStatus string

const (
	CourtCaseStatusVoting    CourtCaseStatus = "voting"
	CourtCaseStatusGuilty    CourtCaseStatus = "guilty"
	CourtCaseStatusNotGuilty CourtCaseStatus = "not_guilty"
	CourtCaseStatusSettled   CourtCaseStatus = "settled"
)

type CourtCase struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	PlaintiffID uuid.UUID       `gorm:"type:uuid;not null;index" json:"plaintiff_id"`
	DefendantID uuid.UUID       `gorm:"type:uuid;not null;index" json:"defendant_id"`
	MatchID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"match_id"`
	Reason      string          `gorm:"type:varchar(100);not null" json:"reason"`
	Status      CourtCaseStatus `gorm:"type:varchar(20);default:'voting';not null" json:"status"`
	ExpiresAt   time.Time       `gorm:"not null" json:"expires_at"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type CourtVoteType string

const (
	CourtVoteGuilty    CourtVoteType = "guilty"
	CourtVoteNotGuilty CourtVoteType = "not_guilty"
)

type CourtVote struct {
	CaseID    uuid.UUID     `gorm:"type:uuid;primaryKey" json:"case_id"`
	JurorID   uuid.UUID     `gorm:"type:uuid;primaryKey" json:"juror_id"`
	Vote      CourtVoteType `gorm:"type:varchar(20);not null" json:"vote"`
	CreatedAt time.Time     `json:"created_at"`
}

func (c *CourtCase) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
