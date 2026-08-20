// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"strings"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExRatingRepository interface {
	Create(ctx context.Context, rating *models.ExRating) error
	HasRated(ctx context.Context, matchID, targetUserID uuid.UUID) (bool, error)
	GetSummaryByUserID(ctx context.Context, targetUserID uuid.UUID) (float64, int64, map[string]int, error)
}

type exRatingRepository struct {
	db *gorm.DB
}

func NewExRatingRepository(db *gorm.DB) ExRatingRepository {
	return &exRatingRepository{db: db}
}

func (r *exRatingRepository) Create(ctx context.Context, rating *models.ExRating) error {
	return r.db.WithContext(ctx).Create(rating).Error
}

func (r *exRatingRepository) HasRated(ctx context.Context, matchID, targetUserID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ExRating{}).
		Where("match_id = ? AND target_user_id = ?", matchID, targetUserID).
		Count(&count).Error
	return count > 0, err
}

func (r *exRatingRepository) GetSummaryByUserID(ctx context.Context, targetUserID uuid.UUID) (float64, int64, map[string]int, error) {
	var ratings []models.ExRating
	err := r.db.WithContext(ctx).Where("target_user_id = ?", targetUserID).Find(&ratings).Error
	if err != nil {
		return 0, 0, nil, err
	}

	total := int64(len(ratings))
	if total == 0 {
		return 0, 0, map[string]int{}, nil
	}

	var sum int
	tagCounts := make(map[string]int)

	for _, rating := range ratings {
		sum += rating.RatingScore
		tags := rating.GetTags()
		for _, tag := range tags {
			if tag != "" {
				tagCounts[tag]++
			}
		}
	}

	avg := float64(sum) / float64(total)
	return avg, total, tagCounts, nil
}
