package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPremiumRepository interface {
	IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error)
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
