// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type LLMClient interface {
	GenerateWingmanSuggestions(ctx context.Context, history []string) ([]Suggestion, error)
}

type Suggestion struct {
	Tone string `json:"tone"`
	Text string `json:"text"`
}

type openAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey string) LLMClient {
	return &openAIClient{
		apiKey: apiKey,
		model:  "gpt-4o-mini", // Using the fast and cheap model
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	ResponseFmt *respFmt  `json:"response_format,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type respFmt struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *openAIClient) GenerateWingmanSuggestions(ctx context.Context, history []string) ([]Suggestion, error) {
	if c.apiKey == "" {
		// Mock implementation for development/testing if API key is not set
		return []Suggestion{
			{Tone: "Hài hước", Text: "Bạn ăn cơm chưa? Mình thì ăn trọn nỗi nhớ bạn rồi 😂"},
			{Tone: "Thả thính", Text: "Trời thu se lạnh, bạn có lạnh không để mình sưởi ấm?"},
			{Tone: "Thẳng thắn", Text: "Nói thật nhé, mình thấy bạn rất thú vị. Bạn nghĩ sao về buổi hẹn tối nay?"},
		}, nil
	}

	prompt := "Dựa vào đoạn lịch sử hội thoại sau (nếu có), hãy đóng vai một trợ lý tình yêu cá tính 'Mỏ Hỗn'. " +
		"Hãy gợi ý 3 tin nhắn phản hồi theo 3 phong cách: 'Hài hước', 'Thả thính', 'Thẳng thắn'. " +
		"Phải trả về dữ liệu chuẩn JSON định dạng array các object có 2 key: 'tone' và 'text'. " +
		"Ví dụ: [{\"tone\": \"Hài hước\", \"text\": \"haha\"}]. Không chứa thêm văn bản nào khác ngoài JSON.\n\n"
	
	if len(history) > 0 {
		prompt += "Lịch sử chat:\n"
		for _, msg := range history {
			prompt += "- " + msg + "\n"
		}
	} else {
		prompt += "Chưa có lịch sử chat nào, hãy gợi ý câu mở lời (ice-breaker)!"
	}

	reqBody := chatRequest{
		Model: c.model,
		Messages: []message{
			{Role: "system", Content: prompt},
		},
		ResponseFmt: &respFmt{Type: "json_object"}, // Force JSON response if supported
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status: %d", resp.StatusCode)
	}

	var resData chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return nil, err
	}

	if len(resData.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}

	content := resData.Choices[0].Message.Content
	
	// Try to unmarshal into []Suggestion. Note: json_object forces an object, but we asked for an array inside. 
	// We might need a wrapper if we used JSON mode. Let's try direct array unmarshal first.
	// Actually, if we use response_format={"type":"json_object"}, it MUST be an object.
	// So let's wrap it in an object format.
	var parsed struct {
		Suggestions []Suggestion `json:"suggestions"`
	}
	
	// In case the model returns the array directly despite json_object (which normally throws an error if it's not an object)
	// We'll try to parse array first, then object.
	var directArray []Suggestion
	if err := json.Unmarshal([]byte(content), &directArray); err == nil {
		return directArray, nil
	}

	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return parsed.Suggestions, nil
}
