// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(ctx context.Context, match *models.Match) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error)
	FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error)
	FindByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Match, error)
	UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ResetStreakForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error
	ResetIslandLevelForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(ctx context.Context, match *models.Match) error {
	return GetDB(ctx, r.db).Create(match).Error
}

func (r *matchRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	var match models.Match
	err := GetDB(ctx, r.db).First(&match, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	var match models.Match
	err := GetDB(ctx, r.db).Unscoped().First(&match, "id = ?", id).Error
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

func (r *matchRepository) FindByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Match, error) {
	var match models.Match
	err := GetDB(ctx, r.db).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1ID, user2ID, user2ID, user1ID).
		First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).
		Where("id = ?", id).
		Delete(&models.Match{}).Error
}

func (r *matchRepository) ResetStreakForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error {
	cutoffTime := time.Now().Add(-inactiveDuration)
	return GetDB(ctx, r.db).
		Model(&models.Match{}).
		Where("last_interaction_at < ?", cutoffTime).
		Where("streak_score > 0").
		Update("streak_score", 0).Error
}

func (r *matchRepository) ResetIslandLevelForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error {
	cutoffTime := time.Now().Add(-inactiveDuration)
	return GetDB(ctx, r.db).
		Model(&models.Match{}).
		Where("last_interaction_at < ?", cutoffTime).
		Where("island_level > 1").
		Update("island_level", 1).Error
}

