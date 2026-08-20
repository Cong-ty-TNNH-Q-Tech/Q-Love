// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPremiumRepository interface {
	IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error)
	ActivatePremium(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error
}

type userPremiumRepository struct {
	db *gorm.DB
}

func NewUserPremiumRepository(db *gorm.DB) UserPremiumRepository {
	return &userPremiumRepository{db: db}
}

func (r *userPremiumRepository) IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := GetDB(ctx, r.db).
		Table("user_premiums").
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userPremiumRepository) ActivatePremium(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error {
	db := GetDB(ctx, r.db)
	
	// Create or update user premium
	// In GORM, we can use Clauses(clause.OnConflict{}) but let's do a simple find/update or create
	// Since we don't have models imported, wait, I need to import models!
	// Wait, is there a UserPremium model?
	// Let's use Raw SQL to avoid importing models if we don't know the exact struct.
	// Actually we should just execute the update or insert.
	
	err := db.WithContext(ctx).Exec(`
		INSERT INTO user_premiums (user_id, expires_at) 
		VALUES (?, ?) 
		ON CONFLICT (user_id) 
		DO UPDATE SET expires_at = CASE 
            WHEN user_premiums.expires_at > CURRENT_TIMESTAMP THEN user_premiums.expires_at + (EXCLUDED.expires_at - CURRENT_TIMESTAMP)
            ELSE EXCLUDED.expires_at
        END`, 

		userID, expiresAt).Error
		
	return err
}
