// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package repository

import (
	"context"

	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type ChatLockRepository interface {
	Create(ctx context.Context, lock *models.ChatLock) error
}

type chatLockRepository struct {
	db *gorm.DB
}

func NewChatLockRepository(db *gorm.DB) ChatLockRepository {
	return &chatLockRepository{db: db}
}

func (r *chatLockRepository) Create(ctx context.Context, lock *models.ChatLock) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(lock).Error
}
