package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatMessage struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MatchID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"match_id"`
	SenderID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"sender_id"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"` // e.g., "text", "locket"
	Content   string         `gorm:"type:text" json:"content"`
	BlurURL   string         `gorm:"type:text" json:"blur_url"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
