// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"github.com/google/uuid"
)

type Landmark struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name               string    `json:"name" gorm:"type:varchar(100)"`
	RadiusMeters       int       `json:"radius_meters" gorm:"default:200"`
	CurrentOwnerClanID *uuid.UUID `json:"current_owner_clan_id" gorm:"type:uuid"`
}
