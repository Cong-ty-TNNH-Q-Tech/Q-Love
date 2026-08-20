// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"mime/multipart"
	"testing"
)

func TestNSFWService_CheckNSFW_Fallback(t *testing.T) {
	service := NewNSFWService(nil)

	// Test case: Fallback without client
	isNSFW, ratio, err := service.CheckNSFW(context.Background(), &multipart.FileHeader{Filename: "test.jpg"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isNSFW {
		t.Errorf("Expected isNSFW=false for fallback, got true")
	}
	if ratio != 0.10 {
		t.Errorf("Expected ratio=0.10, got %v", ratio)
	}
}
