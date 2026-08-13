package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPremium struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PackageName string         `gorm:"type:varchar(255);not null" json:"package_name"`
	ExpiresAt   time.Time      `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
