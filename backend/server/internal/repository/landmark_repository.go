// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"gorm.io/gorm"
)

type LandmarkRepository interface {
	UpdateAllOwners(ctx context.Context, ownerClanID *models.Clan) error
}

type landmarkRepository struct {
	db *gorm.DB
}

func NewLandmarkRepository(db *gorm.DB) LandmarkRepository {
	return &landmarkRepository{db: db}
}

func (r *landmarkRepository) UpdateAllOwners(ctx context.Context, topClan *models.Clan) error {
	var ownerID *string
	if topClan != nil {
		idStr := topClan.ID.String()
		ownerID = &idStr
	}
	return GetDB(ctx, r.db).WithContext(ctx).Model(&models.Landmark{}).Where("1 = 1").Update("current_owner_clan_id", ownerID).Error
}
