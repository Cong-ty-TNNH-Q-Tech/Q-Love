// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/ai"
)

type AIWingmanService interface {
	GetSuggestions(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error)
}

type aiWingmanService struct {
	chatRepo  repository.ChatRepository
	llmClient ai.LLMClient
}

func NewAIWingmanService(chatRepo repository.ChatRepository, llmClient ai.LLMClient) AIWingmanService {
	return &aiWingmanService{
		chatRepo:  chatRepo,
		llmClient: llmClient,
	}
}

func (s *aiWingmanService) GetSuggestions(ctx context.Context, matchID uuid.UUID) ([]ai.Suggestion, error) {
	// 1. Lấy 10 tin nhắn gần nhất
	messages, err := s.chatRepo.GetRecentMessages(ctx, matchID, 10)
	if err != nil {
		return nil, err
	}

	// 2. Format lịch sử
	var history []string
	for _, msg := range messages {
		// Ở môi trường thực tế, chúng ta sẽ cần join tên người gửi. 
		// Nhưng để đơn giản, ta chỉ gộp nội dung tin nhắn.
		history = append(history, msg.Content)
	}

	// 3. Gọi LLM
	return s.llmClient.GenerateWingmanSuggestions(ctx, history)
}
