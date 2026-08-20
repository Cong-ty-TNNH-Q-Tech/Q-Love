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

type CourtRepository interface {
	CreateCase(ctx context.Context, courtCase *models.CourtCase) error
	GetCaseByID(ctx context.Context, id uuid.UUID) (*models.CourtCase, error)
	GetActiveCases(ctx context.Context, limit int) ([]models.CourtCase, error)
	GetExpiredVotingCases(ctx context.Context) ([]models.CourtCase, error)
	UpdateCaseStatus(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error
	CreateVote(ctx context.Context, vote *models.CourtVote) error
	HasUserVoted(ctx context.Context, caseID uuid.UUID, jurorID uuid.UUID) (bool, error)
	CountVotesByCase(ctx context.Context, caseID uuid.UUID) (int64, int64, error)
}

type courtRepository struct {
	db *gorm.DB
}

func NewCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{db: db}
}

func (r *courtRepository) CreateCase(ctx context.Context, courtCase *models.CourtCase) error {
	return r.db.WithContext(ctx).Create(courtCase).Error
}

func (r *courtRepository) GetCaseByID(ctx context.Context, id uuid.UUID) (*models.CourtCase, error) {
	var courtCase models.CourtCase
	err := r.db.WithContext(ctx).First(&courtCase, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &courtCase, nil
}

func (r *courtRepository) GetActiveCases(ctx context.Context, limit int) ([]models.CourtCase, error) {
	var cases []models.CourtCase
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at > ?", models.CourtCaseStatusVoting, time.Now()).
		Limit(limit).
		Find(&cases).Error
	return cases, err
}

func (r *courtRepository) GetExpiredVotingCases(ctx context.Context) ([]models.CourtCase, error) {
	var cases []models.CourtCase
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at <= ?", models.CourtCaseStatusVoting, time.Now()).
		Find(&cases).Error
	return cases, err
}

func (r *courtRepository) UpdateCaseStatus(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.CourtCase{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *courtRepository) CreateVote(ctx context.Context, vote *models.CourtVote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *courtRepository) HasUserVoted(ctx context.Context, caseID uuid.UUID, jurorID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CourtVote{}).
		Where("case_id = ? AND juror_id = ?", caseID, jurorID).
		Count(&count).Error
	return count > 0, err
}

func (r *courtRepository) CountVotesByCase(ctx context.Context, caseID uuid.UUID) (int64, int64, error) {
	var total int64
	var guilty int64

	err := r.db.WithContext(ctx).Model(&models.CourtVote{}).Where("case_id = ?", caseID).Count(&total).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.WithContext(ctx).Model(&models.CourtVote{}).
		Where("case_id = ? AND vote = ?", caseID, models.CourtVoteGuilty).
		Count(&guilty).Error

	return total, guilty, err
}
