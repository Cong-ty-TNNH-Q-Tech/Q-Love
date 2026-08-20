// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Type        string    `json:"type" gorm:"type:varchar(50);not null"`
	Payload     string    `json:"payload" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'sent'"`
	ReferenceID *uuid.UUID `json:"reference_id" gorm:"type:uuid"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
}
