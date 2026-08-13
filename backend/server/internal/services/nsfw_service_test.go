// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"mime/multipart"
	"testing"
)

func TestNSFWService_CheckNSFW(t *testing.T) {
	service := NewNSFWService()

	// Test case: NSFW file
	isNSFW, ratio, err := service.CheckNSFW(context.Background(), &multipart.FileHeader{Filename: "test_nsfw.jpg"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !isNSFW {
		t.Errorf("Expected isNSFW=true for nsfw file, got false")
	}
	if ratio != 0.45 {
		t.Errorf("Expected ratio=0.45, got %v", ratio)
	}

	// Test case: Normal file
	isNSFW, ratio, err = service.CheckNSFW(context.Background(), &multipart.FileHeader{Filename: "normal.jpg"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isNSFW {
		t.Errorf("Expected isNSFW=false for normal file, got true")
	}
	if ratio != 0.10 {
		t.Errorf("Expected ratio=0.10, got %v", ratio)
	}
}
