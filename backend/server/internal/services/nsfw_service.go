// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"io"
	"mime/multipart"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	appconfig "github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
)

type NSFWService interface {
	CheckNSFW(ctx context.Context, file *multipart.FileHeader) (isNSFW bool, skinRatio float64, err error)
}

type RekognitionAPI interface {
	DetectModerationLabels(ctx context.Context, params *rekognition.DetectModerationLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectModerationLabelsOutput, error)
}

type nsfwService struct {
	client RekognitionAPI
}

func NewNSFWService(cfg *appconfig.Config) NSFWService {
	if cfg != nil && cfg.AWSAccessKeyID != "" {
		awsCfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.AWSRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")),
		)
		if err == nil {
			client := rekognition.NewFromConfig(awsCfg)
			return &nsfwService{
				client: client,
			}
		}
	}

	return &nsfwService{
		client: nil,
	}
}

func (s *nsfwService) CheckNSFW(ctx context.Context, file *multipart.FileHeader) (bool, float64, error) {
	if s.client == nil {
		// Fallback to mock for local testing if no AWS credentials are provided
		return false, 0.10, nil
	}

	if file.Size > 5*1024*1024 {
		return false, 0, fmt.Errorf("file too large, max size is 5MB")
	}

	src, err := file.Open()
	if err != nil {
		return false, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return false, 0, fmt.Errorf("failed to read file: %w", err)
	}

	resp, err := s.client.DetectModerationLabels(ctx, &rekognition.DetectModerationLabelsInput{
		Image:         &types.Image{Bytes: data},
		MinConfidence: aws.Float32(70),
	})
	if err != nil {
		return false, 0, fmt.Errorf("failed to detect moderation labels: %w", err)
	}

	for _, label := range resp.ModerationLabels {
		// AWS Rekognition uses top-level categories like "Explicit Nudity", "Suggestive", "Violence"
		name := ""
		if label.Name != nil {
			name = *label.Name
		}
		if name == "Explicit Nudity" || name == "Suggestive" || name == "Adult" {
			confidence := float64(0)
			if label.Confidence != nil {
				confidence = float64(*label.Confidence) / 100.0
			}
			return true, confidence, nil
		}
	}

	return false, 0, nil
}
