// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"encoding/json"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type mockPresigner struct {
	shouldError bool
}

func (m *mockPresigner) GeneratePresignedURL(ctx context.Context, objectKey string, contentType string, contentLength int64) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	return "http://mock.url", nil
}

// We won't test the actual R2 Client upload here (that's for integration),
// but we will test the validation logic of the handler.

func TestGenerateUploadURL_Validation(t *testing.T) {
	app := fiber.New()
	
	// Create mock client
	mockClient := &mockPresigner{}
	h := NewUploadHandler(mockClient)
	app.Post("/upload", h.GenerateUploadURL)

	tests := []struct {
		name          string
		payload       PresignedURLRequest
		expectedCode  int
		expectedError string
	}{

		{
			name: "Valid request",
			payload: PresignedURLRequest{
				Filename:      "test.jpg",
				ContentType:   "image/jpeg",
				ContentLength: 5000000,
			},
			expectedCode: 200,
		},
		{
			name: "Invalid content type",
			payload: PresignedURLRequest{
				Filename:      "test.pdf",
				ContentType:   "application/pdf",
				ContentLength: 5000000,
			},
			expectedCode:  400,
			expectedError: "Only JPG and PNG formats are allowed",
		},
		{
			name: "File size too large",
			payload: PresignedURLRequest{
				Filename:      "test.jpg",
				ContentType:   "image/jpeg",
				ContentLength: 15 * 1024 * 1024, // 15MB
			},
			expectedCode:  400,
			expectedError: "File size must be between 1 byte and 10MB",
		},
		{
			name: "File size too small (0)",
			payload: PresignedURLRequest{
				Filename:      "test.jpg",
				ContentType:   "image/jpeg",
				ContentLength: 0,
			},
			expectedCode:  400,
			expectedError: "File size must be between 1 byte and 10MB",
		},
		{
			name: "No extension fallback",
			payload: PresignedURLRequest{
				Filename:      "test",
				ContentType:   "image/jpeg",
				ContentLength: 5000000,
			},
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/upload", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}



			if tt.name == "Valid request" {
				if resp.StatusCode != 200 {
					t.Errorf("Expected status 200, got %d", resp.StatusCode)
				}
				return
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			var respBody map[string]string
			json.NewDecoder(resp.Body).Decode(&respBody)

			if respBody["error"] != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, respBody["error"])
			}
		})
	}

	// Test Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload", bytes.NewReader([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400 for invalid JSON, got %d", resp.StatusCode)
		}
	})

	// Test Internal Error (Mock Error)
	t.Run("Presigner Error", func(t *testing.T) {
		appErr := fiber.New()
		errClient := &mockPresigner{shouldError: true}
		hErr := NewUploadHandler(errClient)
		appErr.Post("/upload", hErr.GenerateUploadURL)

		payload := PresignedURLRequest{
			Filename:      "test.jpg",
			ContentType:   "image/jpeg",
			ContentLength: 5000000,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/upload", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, _ := appErr.Test(req, -1)
		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500 for presigner error, got %d", resp.StatusCode)
		}
	})
}
