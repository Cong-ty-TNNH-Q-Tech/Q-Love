// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type AIWingmanService interface {
	SuggestReplies(ctx context.Context, matchID uuid.UUID) ([]string, error)
	MaskPII(text string) string
}

type aiWingmanService struct {
	chatRepo   repository.ChatRepository // Notice this uses ChatRepository!
	openAPIKey string
	httpClient *http.Client
}

func NewAIWingmanService(chatRepo repository.ChatRepository, openAPIKey string) AIWingmanService {
	return &aiWingmanService{
		chatRepo:   chatRepo,
		openAPIKey: openAPIKey,
		httpClient: &http.Client{},
	}
}

// MaskPII removes phone numbers and emails to protect privacy
func (s *aiWingmanService) MaskPII(text string) string {
	// Simple email regex
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	text = emailRegex.ReplaceAllString(text, "[HIDDEN_EMAIL]")

	// Simple phone regex (10-11 digits)
	phoneRegex := regexp.MustCompile(`(?:\+84|0)(?:\d[\s\-\.]?){8,10}\d`)
	text = phoneRegex.ReplaceAllString(text, "[HIDDEN_PHONE]")

	return text
}

func (s *aiWingmanService) SuggestReplies(ctx context.Context, matchID uuid.UUID) ([]string, error) {
	if s.openAPIKey == "" {
		// Mock logic for local testing without API Key
		return []string{
			"Gợi ý 1 (Mock): Chào em, em ăn cơm chưa?",
			"Gợi ý 2 (Mock): Sao nhắn tin chậm thế, bận à?",
			"Gợi ý 3 (Mock): Đi chơi không em?",
		}, nil
	}

	messages, err := s.chatRepo.GetMessagesByMatchID(ctx, matchID, 5, nil)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, errors.New("không có tin nhắn nào để gợi ý")
	}

	// Format conversation context
	var conversationBuilder strings.Builder
	for i := len(messages) - 1; i >= 0; i-- { // Reverse to chronological order
		msg := messages[i]
		sender := "Người A" // Need actual logic to distinguish user, but let's just use generic labels for context
		content := s.MaskPII(msg.Content)
		conversationBuilder.WriteString(fmt.Sprintf("%s: %s\n", sender, content))
	}

	prompt := fmt.Sprintf(`Đọc đoạn hội thoại sau và đóng vai "Trợ lý Mỏ Hỗn" (một AI vui tính, hơi xéo xắt, hài hước, đôi khi thả thính). Hãy gợi ý 3 câu trả lời tiếp theo thật ngắn gọn, mỗi câu trên một dòng, cách nhau bởi ký tự newline (\n). Không giải thích thêm.
Hội thoại:
%s
Gợi ý:`, conversationBuilder.String())

	reqBody := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "Bạn là trợ lý AI chuyên gợi ý tin nhắn hẹn hò."},
			{"role": "user", "content": prompt},
		},
		"max_tokens": 150,
	}

	jsonValue, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.openAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) > 0 {
		content := result.Choices[0].Message.Content
		replies := strings.Split(content, "\n")
		var cleanReplies []string
		for _, r := range replies {
			r = strings.TrimSpace(r)
			if r != "" {
				// Remove numbering like "1. ", "2. ", "- "
				r = regexp.MustCompile(`^(\d+\. |- )`).ReplaceAllString(r, "")
				cleanReplies = append(cleanReplies, r)
			}
		}
		if len(cleanReplies) > 3 {
			cleanReplies = cleanReplies[:3]
		}
		return cleanReplies, nil
	}

	return nil, errors.New("không có kết quả trả về từ OpenAI")
}
