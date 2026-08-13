// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/ai"
)

type mockChatRepository struct {
	getRecentMessagesFn func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error)
}

func (m *mockChatRepository) GetRecentMessages(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	return m.getRecentMessagesFn(ctx, matchID, limit)
}

type mockLLMClient struct {
	generateFn func(ctx context.Context, history []string) ([]ai.Suggestion, error)
}

func (m *mockLLMClient) GenerateWingmanSuggestions(ctx context.Context, history []string) ([]ai.Suggestion, error) {
	return m.generateFn(ctx, history)
}

func TestAIWingmanService_GetSuggestions(t *testing.T) {
	mockChat := &mockChatRepository{
		getRecentMessagesFn: func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error) {
			return []models.ChatMessage{
				{Content: "Hi"},
				{Content: "Hello"},
			}, nil
		},
	}
	mockLLM := &mockLLMClient{
		generateFn: func(ctx context.Context, history []string) ([]ai.Suggestion, error) {
			return []ai.Suggestion{
				{Tone: "Hài hước", Text: "Haha"},
			}, nil
		},
	}

	svc := NewAIWingmanService(mockChat, mockLLM)
	suggestions, err := svc.GetSuggestions(context.Background(), uuid.New())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestAIWingmanService_GetSuggestions_ChatError(t *testing.T) {
	mockChat := &mockChatRepository{
		getRecentMessagesFn: func(ctx context.Context, matchID uuid.UUID, limit int) ([]models.ChatMessage, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewAIWingmanService(mockChat, &mockLLMClient{})
	
	_, err := svc.GetSuggestions(context.Background(), uuid.New())
	if err == nil {
		t.Error("Expected error")
	}
}
