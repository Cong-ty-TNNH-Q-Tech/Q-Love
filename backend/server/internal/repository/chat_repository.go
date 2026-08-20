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

type ChatRepository interface {
	Create(ctx context.Context, msg *models.ChatMessage) error
	GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error)
	CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(ctx context.Context, msg *models.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatRepository) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	
	query := r.db.WithContext(ctx).Where("match_id = ?", matchID)
	
	if before != nil {
		query = query.Where("created_at < ?", *before)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *chatRepository) CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ChatMessage{}).Where("match_id = ?", matchID).Count(&count).Error
	return count, err
}
