// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"mime/multipart"
	"strings"
)

type NSFWService interface {
	CheckNSFW(ctx context.Context, file *multipart.FileHeader) (isNSFW bool, skinRatio float64, err error)
}

type nsfwService struct{}

func NewNSFWService() NSFWService {
	return &nsfwService{}
}

func (s *nsfwService) CheckNSFW(ctx context.Context, file *multipart.FileHeader) (bool, float64, error) {
	// Mock Implementation for NSFW AI Check
	// In a real environment, this would send the file to an AI service (e.g. AWS Rekognition, or internal Python gRPC)
	// For testing, we mock it: if filename contains "nsfw", we flag it with a mock high skin ratio
	if strings.Contains(strings.ToLower(file.Filename), "nsfw") {
		return true, 0.45, nil // 45% skin ratio > 30%
	}
	return false, 0.10, nil // 10% skin ratio <= 30%
}
