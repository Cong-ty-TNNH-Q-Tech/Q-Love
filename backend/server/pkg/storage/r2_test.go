package storage

import (
	"context"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
)

func TestNewR2Client(t *testing.T) {
	cfg := &config.Config{
		R2AccountID:       "dummy",
		R2AccessKeyID:     "dummy",
		R2SecretAccessKey: "dummy",
		R2BucketName:      "dummy",
	}
	client, err := NewR2Client(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	url, err := client.GeneratePresignedURL(context.Background(), "test.jpg", "image/jpeg", 100)
	if err != nil {
		t.Fatalf("Expected no error for presigning offline, got %v", err)
	}
	if url == "" {
		t.Error("Expected a url, got empty string")
	}
}
