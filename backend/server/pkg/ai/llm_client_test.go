// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestOpenAIClient_Mock(t *testing.T) {
	client := NewOpenAIClient("") // Empty API key should return mock data
	suggestions, err := client.GenerateWingmanSuggestions(context.Background(), []string{"Hello"})
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(suggestions) != 3 {
		t.Fatalf("Expected 3 suggestions, got %d", len(suggestions))
	}
	
	if suggestions[0].Tone != "Hài hước" {
		t.Errorf("Expected first tone to be 'Hài hước', got %s", suggestions[0].Tone)
	}
}

func TestOpenAIClient_Success(t *testing.T) {
	client := NewOpenAIClient("test-key").(*openAIClient)
	
	// Mock HTTP Client
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		mockResp := chatResponse{}
		mockResp.Choices = append(mockResp.Choices, struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			Message: struct {
				Content string `json:"content"`
			}{
				Content: `{"suggestions": [{"tone": "Hài hước", "text": "Hello"}]}`,
			},
		})
		
		b, _ := json.Marshal(mockResp)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	suggestions, err := client.GenerateWingmanSuggestions(context.Background(), []string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("Expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestOpenAIClient_DirectArraySuccess(t *testing.T) {
	client := NewOpenAIClient("test-key").(*openAIClient)
	
	// Mock HTTP Client
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		mockResp := chatResponse{}
		mockResp.Choices = append(mockResp.Choices, struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			Message: struct {
				Content string `json:"content"`
			}{
				Content: `[{"tone": "Thẳng thắn", "text": "Hi"}]`,
			},
		})
		
		b, _ := json.Marshal(mockResp)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	suggestions, err := client.GenerateWingmanSuggestions(context.Background(), []string{"History 1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(suggestions) != 1 {
		t.Fatalf("Expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestOpenAIClient_ErrorStatus(t *testing.T) {
	client := NewOpenAIClient("test-key").(*openAIClient)
	
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			Header:     make(http.Header),
		}
	})

	_, err := client.GenerateWingmanSuggestions(context.Background(), []string{"Hello"})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestOpenAIClient_EmptyChoices(t *testing.T) {
	client := NewOpenAIClient("test-key").(*openAIClient)
	
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		mockResp := chatResponse{}
		b, _ := json.Marshal(mockResp)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	_, err := client.GenerateWingmanSuggestions(context.Background(), []string{"Hello"})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestOpenAIClient_ParseError(t *testing.T) {
	client := NewOpenAIClient("test-key").(*openAIClient)
	
	client.httpClient = NewTestClient(func(req *http.Request) *http.Response {
		mockResp := chatResponse{}
		mockResp.Choices = append(mockResp.Choices, struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			Message: struct {
				Content string `json:"content"`
			}{
				Content: `invalid json`,
			},
		})
		
		b, _ := json.Marshal(mockResp)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	_, err := client.GenerateWingmanSuggestions(context.Background(), []string{"Hello"})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
