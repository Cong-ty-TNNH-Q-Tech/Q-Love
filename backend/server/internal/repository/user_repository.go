// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"
	"errors"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error)
	Create(ctx context.Context, user *models.User) error
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
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

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := GetDB(ctx, r.db).WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found instead of error
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := GetDB(ctx, r.db).WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
