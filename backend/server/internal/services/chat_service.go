// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.
package services

import (
	"context"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidMessage = errors.New("invalid chat message")
)

type ChatService interface {
	SaveMessage(ctx context.Context, senderID, matchID uuid.UUID, msgType, content string) (*models.ChatMessage, error)
	GetMessages(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) SaveMessage(ctx context.Context, senderID, matchID uuid.UUID, msgType, content string) (*models.ChatMessage, error) {
	if content == "" && msgType == "text" {
		return nil, ErrInvalidMessage
	}

	msg := &models.ChatMessage{
		ID:       uuid.New(),
		MatchID:  matchID,
		SenderID: senderID,
		Type:     msgType,
		Content:  content,
	}

	if err := s.chatRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.chatRepo.GetMessagesByMatchID(ctx, matchID, limit, before)
}

