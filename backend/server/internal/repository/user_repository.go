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

type UserRepository interface {
	GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error)
	GetSpiritualFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID
	
	// Assuming card_profiles stores metrics
	err := GetDB(ctx, r.db).WithContext(ctx).
		Table("card_profiles").
		Order("(match_count_cached + locket_count_cached + clan_upvote_cached - court_penalty_cached) DESC").
		Limit(limit).
		Pluck("user_id", &userIDs).Error
		
	return userIDs, err
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := GetDB(ctx, r.db).WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) {
	// Find user's location first
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var feed []models.User
	
	err = GetDB(ctx, r.db).WithContext(ctx).
		Where("id != ?", userID).
		Where("is_shadowbanned = ?", false).
		Where("ST_DWithin(location::geography, ?::geography, ?)", user.Location, radius).
		Limit(50).
		Find(&feed).Error
		
	return feed, err
}

func (r *userRepository) GetSpiritualFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var feed []models.User
	err = GetDB(ctx, r.db).WithContext(ctx).
		Where("id != ?", userID).
		Where("is_shadowbanned = ?", false).
		Where("ST_DWithin(location::geography, ?::geography, ?)", user.Location, radius).
		Limit(1000).
		Find(&feed).Error
		
	return feed, err
}
