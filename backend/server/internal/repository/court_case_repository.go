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

type CourtCaseRepository interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type courtCaseRepository struct {
	db *gorm.DB
}

func NewCourtCaseRepository(db *gorm.DB) CourtCaseRepository {
	return &courtCaseRepository{db: db}
}

func (r *courtCaseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return GetDB(ctx, r.db).Model(&models.CourtCase{}).Where("id = ?", id).Update("status", status).Error
}
