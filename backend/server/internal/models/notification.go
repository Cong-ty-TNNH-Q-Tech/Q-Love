// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string     `gorm:"type:varchar(50);not null" json:"type"`
	Payload     string     `gorm:"type:text" json:"payload"`
	Status      string     `gorm:"type:varchar(20);default:'sent'" json:"status"`
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
