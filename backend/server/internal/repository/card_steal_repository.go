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

type CardStealRepository interface {
	Create(ctx context.Context, steal *models.CardSteal) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.CardSteal, error)
	UpdateResult(ctx context.Context, id uuid.UUID, result string) error
	TransferCardOwnership(ctx context.Context, collectorID uuid.UUID, targetCardID uuid.UUID) error
}

type cardStealRepository struct {
	db *gorm.DB
}

func NewCardStealRepository(db *gorm.DB) CardStealRepository {
	return &cardStealRepository{db: db}
}

func (r *cardStealRepository) Create(ctx context.Context, steal *models.CardSteal) error {
	return GetDB(ctx, r.db).Create(steal).Error
}

func (r *cardStealRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.CardSteal, error) {
	var steal models.CardSteal
	err := GetDB(ctx, r.db).Where("id = ?", id).First(&steal).Error
	if err != nil {
		return nil, err
	}
	return &steal, nil
}

func (r *cardStealRepository) UpdateResult(ctx context.Context, id uuid.UUID, result string) error {
	return GetDB(ctx, r.db).Model(&models.CardSteal{}).Where("id = ?", id).Update("result", result).Error
}

func (r *cardStealRepository) TransferCardOwnership(ctx context.Context, collectorID uuid.UUID, targetCardID uuid.UUID) error {
	// Transfer ownership by recording a card_transactions
	tx := &models.CardTransaction{
		CollectorID:        collectorID,
		TargetUserID:       targetCardID,
		Type:               "steal",
		Quantity:           1,
		PriceAtTransaction: 0.0,
	}
	return GetDB(ctx, r.db).Create(tx).Error
}
