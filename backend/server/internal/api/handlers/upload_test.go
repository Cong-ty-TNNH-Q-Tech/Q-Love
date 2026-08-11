package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

// We won't test the actual R2 Client upload here (that's for integration),
// but we will test the validation logic of the handler.

func TestGenerateUploadURL_Validation(t *testing.T) {
	app := fiber.New()
	
	// Create dummy client to avoid panic
	cfg := &config.Config{
		R2AccountID:       "dummy",
		R2AccessKeyID:     "dummy",
		R2SecretAccessKey: "dummy",
		R2BucketName:      "dummy",
	}
	dummyClient, _ := storage.NewR2Client(cfg)
	h := NewUploadHandler(dummyClient)
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
}
