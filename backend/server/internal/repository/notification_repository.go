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

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return GetDB(ctx, r.db).WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("status", status).Error
}
