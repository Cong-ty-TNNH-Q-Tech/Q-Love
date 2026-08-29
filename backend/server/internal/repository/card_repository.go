// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardRepository interface {
	GetCardProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error)
	GetCardProfileForUpdate(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error)
	CreateCardProfile(ctx context.Context, profile *models.CardProfile) error
	UpdateCardProfile(ctx context.Context, profile *models.CardProfile) error
	CreateCardTransaction(ctx context.Context, tx *models.CardTransaction) error
	GetOwnedQuantity(ctx context.Context, collectorID, targetUserID uuid.UUID) (int, error)
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

// Extract DB from context if inside a transaction
func (r *cardRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *cardRepository) GetCardProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	var profile models.CardProfile
	err := r.getDB(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *cardRepository) GetCardProfileForUpdate(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	var profile models.CardProfile
	err := r.getDB(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *cardRepository) CreateCardProfile(ctx context.Context, profile *models.CardProfile) error {
	return r.getDB(ctx).Create(profile).Error
}

func (r *cardRepository) UpdateCardProfile(ctx context.Context, profile *models.CardProfile) error {
	return r.getDB(ctx).Save(profile).Error
}

func (r *cardRepository) CreateCardTransaction(ctx context.Context, transaction *models.CardTransaction) error {
	return r.getDB(ctx).Create(transaction).Error
}

func (r *cardRepository) GetOwnedQuantity(ctx context.Context, collectorID, targetUserID uuid.UUID) (int, error) {
	var buyQty int
	r.getDB(ctx).Model(&models.CardTransaction{}).Where("collector_id = ? AND target_user_id = ? AND type = ?", collectorID, targetUserID, "buy").Select("COALESCE(SUM(quantity), 0)").Scan(&buyQty)
	
	var sellQty int
	r.getDB(ctx).Model(&models.CardTransaction{}).Where("collector_id = ? AND target_user_id = ? AND type = ?", collectorID, targetUserID, "sell").Select("COALESCE(SUM(quantity), 0)").Scan(&sellQty)
	
	return buyQty - sellQty, nil
}
