// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Phone          string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Name           string         `gorm:"type:varchar(100)" json:"name"`
	DOB            *time.Time     `gorm:"type:date" json:"dob,omitempty"`
	Gender         string         `gorm:"type:varchar(20)" json:"gender,omitempty"`
	Location       string         `gorm:"type:geometry(Point,4326)" json:"-"` // We ignore json export for raw WKB geometry
	Level          int            `gorm:"default:1" json:"level"`
	IsShadowbanned bool           `gorm:"default:false" json:"is_shadowbanned"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
