// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
)

type R2Client struct {
	S3Client      *s3.Client
	BucketName    string
	PresignClient *s3.PresignClient
}

func NewR2Client(cfg *appConfig.Config) (*R2Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(client)

	return &R2Client{
		S3Client:      client,
		BucketName:    cfg.R2BucketName,
		PresignClient: presignClient,
	}, nil
}

func (r *R2Client) GeneratePresignedURL(ctx context.Context, objectKey string, contentType string, contentLength int64) (string, error) {
	request, err := r.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.BucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * time.Minute)
	})
	if err != nil {
		return "", fmt.Errorf("couldn't get a presigned request: %v", err)
	}
	return request.URL, nil
}
