// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type WingmanRepository interface {
	CreateReferral(ctx context.Context, referral *models.WingmanReferral) error
	GetReferralByID(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error)
	UpdateReferral(ctx context.Context, referral *models.WingmanReferral) error
}

type wingmanRepository struct {
	db *gorm.DB
}

func NewWingmanRepository(db *gorm.DB) WingmanRepository {
	return &wingmanRepository{db: db}
}

func (r *wingmanRepository) CreateReferral(ctx context.Context, referral *models.WingmanReferral) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(referral).Error
}

func (r *wingmanRepository) GetReferralByID(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
	var referral models.WingmanReferral
	if err := GetDB(ctx, r.db).WithContext(ctx).First(&referral, "id = ?", referralID).Error; err != nil {
		return nil, err
	}
	return &referral, nil
}

func (r *wingmanRepository) UpdateReferral(ctx context.Context, referral *models.WingmanReferral) error {
	return GetDB(ctx, r.db).WithContext(ctx).Save(referral).Error
}
