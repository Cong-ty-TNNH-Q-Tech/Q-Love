package models

import (
	"time"

	"github.com/google/uuid"
)

type WingmanReferral struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WingmanID  uuid.UUID `gorm:"type:uuid;not null" json:"wingman_id"`
	Target1ID  uuid.UUID `gorm:"type:uuid;not null" json:"target1_id"`
	Target2ID  uuid.UUID `gorm:"type:uuid;not null" json:"target2_id"`
	Status     string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, matched, dated, rewarded
	DeepLink   string    `gorm:"type:varchar(255)" json:"deep_link"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}
