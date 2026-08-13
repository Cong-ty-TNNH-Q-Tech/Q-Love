package repository

import (
	"context"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error)
	UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	var match models.Match
	err := GetDB(ctx, r.db).First(&match, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error {
	return GetDB(ctx, r.db).
		Model(&models.Match{}).
		Where("id = ?", id).
		Update("last_interaction_at", t).Error
}
