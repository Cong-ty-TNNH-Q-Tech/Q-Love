// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error)
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
