// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"time"

	"github.com/google/uuid"
)

type CardTransaction struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CollectorID        uuid.UUID `gorm:"type:uuid;not null;index" json:"collector_id"`
	TargetUserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"target_user_id"`
	Type               string    `gorm:"type:varchar(20)" json:"type"`
	Quantity           int       `json:"quantity"`
	PriceAtTransaction float64   `gorm:"type:numeric(15,2)" json:"price_at_transaction"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}
