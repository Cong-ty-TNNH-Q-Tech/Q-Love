// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"qlove/internal/models"
)

type ShameRepository interface {
	GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShame, error)
	IncrementTomatoCount(ctx context.Context, tx *gorm.DB, shameID uuid.UUID, count int) error
}

type shameRepository struct {
	db *gorm.DB
}

func NewShameRepository(db *gorm.DB) ShameRepository {
	return &shameRepository{db: db}
}

func (r *shameRepository) GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShame, error) {
	var shames []models.WallOfShame
	err := r.db.WithContext(ctx).
		Where("expires_at > ?", time.Now()).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&shames).Error
	return shames, err
}

func (r *shameRepository) IncrementTomatoCount(ctx context.Context, tx *gorm.DB, shameID uuid.UUID, count int) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).
		Model(&models.WallOfShame{}).
		Where("id = ?", shameID).
		UpdateColumn("tomatoes_thrown", gorm.Expr("tomatoes_thrown + ?", count)).Error
}
