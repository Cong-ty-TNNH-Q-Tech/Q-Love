// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
)

type MockChatRepository struct {
	CreateFunc               func(ctx context.Context, msg *models.ChatMessage) error
	GetMessagesByMatchIDFunc func(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error)
	CountMessagesByMatchIDFunc func(ctx context.Context, matchID uuid.UUID) (int64, error)
}

func (m *MockChatRepository) Create(ctx context.Context, msg *models.ChatMessage) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, msg)
	}
	return nil
}

func (m *MockChatRepository) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	if m.GetMessagesByMatchIDFunc != nil {
		return m.GetMessagesByMatchIDFunc(ctx, matchID, limit, before)
	}
	return []models.ChatMessage{}, nil
}

func (m *MockChatRepository) CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error) {
	if m.CountMessagesByMatchIDFunc != nil {
		return m.CountMessagesByMatchIDFunc(ctx, matchID)
	}
	return 0, nil
}

func TestChatService_SaveMessage(t *testing.T) {
	mockRepo := &MockChatRepository{
		CreateFunc: func(ctx context.Context, msg *models.ChatMessage) error {
			return nil
		},
	}
	service := NewChatService(mockRepo)

	senderID := uuid.New()
	matchID := uuid.New()

	msg, err := service.SaveMessage(context.Background(), senderID, matchID, "text", "hello")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if msg == nil {
		t.Errorf("expected message, got nil")
	}

	// Test invalid message
	_, err = service.SaveMessage(context.Background(), senderID, matchID, "text", "")
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("expected ErrInvalidMessage, got %v", err)
	}
}

func TestChatService_GetMessages(t *testing.T) {
	mockRepo := &MockChatRepository{
		GetMessagesByMatchIDFunc: func(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
			return []models.ChatMessage{{ID: uuid.New()}}, nil
		},
	}
	service := NewChatService(mockRepo)

	messages, err := service.GetMessages(context.Background(), uuid.New(), 0, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
}

