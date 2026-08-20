// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ExRating struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TargetUserID uuid.UUID `gorm:"type:uuid;index;not null"`
	MatchID      uuid.UUID `gorm:"type:uuid;index;not null"`
	RatingScore  int       `gorm:"not null"`
	TagsString   string    `gorm:"type:text"` // Comma-separated tags
	CreatedAt    time.Time
}

// GetTags returns the tags as a slice of strings
func (r *ExRating) GetTags() []string {
	if r.TagsString == "" {
		return []string{}
	}
	return strings.Split(r.TagsString, ",")
}

// SetTags stores a slice of strings as a comma-separated string
func (r *ExRating) SetTags(tags []string) {
	r.TagsString = strings.Join(tags, ",")
}
