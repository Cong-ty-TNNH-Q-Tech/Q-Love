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

type DatingContractRepository interface {
	Create(ctx context.Context, contract *models.DatingContract) error
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.DatingContract, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.DatingContract, error)
	Update(ctx context.Context, contract *models.DatingContract) error
}

type datingContractRepository struct {
	db *gorm.DB
}

func NewDatingContractRepository(db *gorm.DB) DatingContractRepository {
	return &datingContractRepository{db: db}
}

func (r *datingContractRepository) Create(ctx context.Context, contract *models.DatingContract) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).Create(contract).Error
}

func (r *datingContractRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.DatingContract, error) {
	db := GetDB(ctx, r.db)
	var contract models.DatingContract
	err := db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *datingContractRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.DatingContract, error) {
	db := GetDB(ctx, r.db)
	var contract models.DatingContract
	err := db.WithContext(ctx).First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *datingContractRepository) Update(ctx context.Context, contract *models.DatingContract) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).Save(contract).Error
}
