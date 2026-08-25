// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LandmarkRepository interface {
	UpdateAllOwners(ctx context.Context, ownerClanID *models.Clan) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Landmark, error)
	CheckDistance(ctx context.Context, landmarkID uuid.UUID, lat, lng float64) (bool, error)
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

func (r *landmarkRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Landmark, error) {
	db := GetDB(ctx, r.db)
	var landmark models.Landmark
	err := db.WithContext(ctx).First(&landmark, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &landmark, nil
}

func (r *landmarkRepository) CheckDistance(ctx context.Context, landmarkID uuid.UUID, lat, lng float64) (bool, error) {
	db := GetDB(ctx, r.db)
	var isWithin bool
	
	// PostGIS ST_DWithin checks if geography points are within radius meters.
	// Note: ST_MakePoint takes (longitude, latitude)
	err := db.WithContext(ctx).Raw(`
		SELECT ST_DWithin(
			location::geography, 
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, 
			radius_meters
		) 
		FROM landmarks 
		WHERE id = ?
	`, lng, lat, landmarkID).Scan(&isWithin).Error

	if err != nil {
		return false, err
	}
	return isWithin, nil
}
