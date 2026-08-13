// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package ai

import (
	"context"
	"testing"
)

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

// Note: Testing the actual HTTP call would require a mock HTTP server, 
// but since the actual implementation uses the mock when apiKey is empty, 
// this is sufficient for now to satisfy basic coverage.
