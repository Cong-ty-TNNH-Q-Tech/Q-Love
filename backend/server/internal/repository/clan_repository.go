// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"gorm.io/gorm"
)

type ClanRepository interface {
	CreateClan(ctx context.Context, clan *models.Clan) error
	AddClanMember(ctx context.Context, member *models.ClanMember) error
	FindByName(ctx context.Context, name string) (*models.Clan, error)
}

type clanRepository struct {
	db *gorm.DB
}

func NewClanRepository(db *gorm.DB) ClanRepository {
	return &clanRepository{db: db}
}

func (r *clanRepository) CreateClan(ctx context.Context, clan *models.Clan) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(clan).Error
}

func (r *clanRepository) AddClanMember(ctx context.Context, member *models.ClanMember) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(member).Error
}

func (r *clanRepository) FindByName(ctx context.Context, name string) (*models.Clan, error) {
	var clan models.Clan
	if err := GetDB(ctx, r.db).WithContext(ctx).Where("name = ?", name).First(&clan).Error; err != nil {
		return nil, err
	}
	return &clan, nil
}
