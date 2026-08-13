package repository

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	Create(ctx context.Context, message *models.ChatMessage) error
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(ctx context.Context, message *models.ChatMessage) error {
	return GetDB(ctx, r.db).Create(message).Error
}
