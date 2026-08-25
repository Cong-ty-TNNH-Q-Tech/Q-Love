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

type UserViolationRepository interface {
	Create(ctx context.Context, violation *models.UserViolation) error
	CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error)
	BanUser(ctx context.Context, userID uuid.UUID) error
	GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error)
	DeleteViolation(ctx context.Context, id uuid.UUID) error
}

type userViolationRepository struct {
	db *gorm.DB
}

func NewUserViolationRepository(db *gorm.DB) UserViolationRepository {
	return &userViolationRepository{db: db}
}

func (r *userViolationRepository) Create(ctx context.Context, violation *models.UserViolation) error {
	return GetDB(ctx, r.db).Create(violation).Error
}

func (r *userViolationRepository) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) {
	var count int64
	err := GetDB(ctx, r.db).Model(&models.UserViolation{}).
		Where("user_id = ? AND type = ? AND is_active = true", userID, vType).
		Count(&count).Error
	return count, err
}

func (r *userViolationRepository) BanUser(ctx context.Context, userID uuid.UUID) error {
	// Execute raw update on users table to enforce auto-ban
	// We use is_shadowbanned per the ERD logic, or you can use is_banned if preferred.
	// Since the requirement states "auto-ban", we will update is_shadowbanned to true.
	return GetDB(ctx, r.db).Table("users").Where("id = ?", userID).Update("is_shadowbanned", true).Error
}

func (r *userViolationRepository) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	var violations []models.UserViolation
	var total int64

	db := GetDB(ctx, r.db).Model(&models.UserViolation{}).Where("is_active = true")
	
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&violations).Error
	return violations, total, err
}

func (r *userViolationRepository) DeleteViolation(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Model(&models.UserViolation{}).Where("id = ?", id).Update("is_active", false).Error
}

