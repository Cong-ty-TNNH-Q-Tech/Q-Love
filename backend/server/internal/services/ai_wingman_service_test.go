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
	"github.com/stretchr/testify/assert"
)

type mockChatRepoForAI struct {
	GetMessagesByMatchIDFn func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error)
}

func (m *mockChatRepoForAI) Create(ctx context.Context, msg *models.ChatMessage) error {
	return nil
}
func (m *mockChatRepoForAI) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	if m.GetMessagesByMatchIDFn != nil {
		return m.GetMessagesByMatchIDFn(ctx, matchID, limit)
	}
	return []models.ChatMessage{
		{Content: "Chào em 0987654321"},
		{Content: "Email anh là test@gmail.com"},
	}, nil
}

func TestAIWingmanService_MaskPII(t *testing.T) {
	svc := NewAIWingmanService(nil, "")
	masked := svc.MaskPII("SĐT của anh là 0912345678, gọi anh nhé hoặc email test@example.com")
	
	assert.Contains(t, masked, "[HIDDEN_PHONE]")
	assert.Contains(t, masked, "[HIDDEN_EMAIL]")
	assert.NotContains(t, masked, "0912345678")
	assert.NotContains(t, masked, "test@example.com")
}

func TestAIWingmanService_SuggestReplies_Mock(t *testing.T) {
	mockRepo := &mockChatRepoForAI{}
	svc := NewAIWingmanService(mockRepo, "") // Empty API Key triggers mock response

	replies, err := svc.SuggestReplies(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Len(t, replies, 3)
	assert.Contains(t, replies[0], "Gợi ý 1 (Mock)")
}

func TestAIWingmanService_SuggestReplies_Error(t *testing.T) {
	mockRepo := &mockChatRepoForAI{
		GetMessagesByMatchIDFn: func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewAIWingmanService(mockRepo, "fake-api-key")

	_, err := svc.SuggestReplies(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestAIWingmanService_SuggestReplies_NoMessages(t *testing.T) {
	mockRepo := &mockChatRepoForAI{
		GetMessagesByMatchIDFn: func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error) {
			return []models.ChatMessage{}, nil
		},
	}
	svc := NewAIWingmanService(mockRepo, "fake-api-key")

	_, err := svc.SuggestReplies(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "không có tin nhắn nào để gợi ý")
}

