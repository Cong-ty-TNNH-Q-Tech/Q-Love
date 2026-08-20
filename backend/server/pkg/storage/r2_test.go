// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package storage

import (
	"context"
	"strings"
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

func TestGeneratePresignedURL_Error(t *testing.T) {
	cfg := &config.Config{
		R2AccountID:       "dummy",
		R2AccessKeyID:     "dummy",
		R2SecretAccessKey: "dummy",
		R2BucketName:      "dummy",
	}
	client, _ := NewR2Client(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately
	
	_, err := client.GeneratePresignedURL(ctx, "test.jpg", "image/jpeg", 100)
	if err == nil {
		t.Fatal("Expected error due to cancelled context, got nil")
	}
}

func TestUploadFile_Error(t *testing.T) {
	cfg := &config.Config{
		R2AccountID:       "dummy",
		R2AccessKeyID:     "dummy",
		R2SecretAccessKey: "dummy",
		R2BucketName:      "dummy",
	}
	client, _ := NewR2Client(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to force network error
	
	_, err := client.UploadFile(ctx, "test.jpg", strings.NewReader("dummy"), "image/jpeg")
	if err == nil {
		t.Fatal("Expected error due to cancelled context, got nil")
	}
}
