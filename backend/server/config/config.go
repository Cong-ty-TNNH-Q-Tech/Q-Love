// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	R2AccountID             string
	R2AccessKeyID           string
	R2SecretAccessKey       string
	R2BucketName            string
	Port                    string
	SentryDSN               string
	Environment             string
	DatabaseDSN             string
	RedisURL                string
	RevenueCatWebhookSecret string
	AWSRegion               string
	AWSAccessKeyID          string
	AWSSecretAccessKey      string
	JWTSecret               string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	dbDSN := os.Getenv("DATABASE_DSN")
	if dbDSN == "" {
		log.Panic("DATABASE_DSN is required")
	}

	rcSecret := os.Getenv("REVENUECAT_WEBHOOK_SECRET")
	if rcSecret == "" {
		log.Panic("REVENUECAT_WEBHOOK_SECRET is required")
	}

	r2AccessKey := os.Getenv("R2_ACCESS_KEY_ID")
	if r2AccessKey == "" {
		log.Panic("R2_ACCESS_KEY_ID is required")
	}

	r2SecretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	if r2SecretKey == "" {
		log.Panic("R2_SECRET_ACCESS_KEY is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Panic("JWT_SECRET is required")
	}

	return &Config{
		R2AccountID:             getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:           r2AccessKey,
		R2SecretAccessKey:       r2SecretKey,
		R2BucketName:            getEnv("R2_BUCKET_NAME", ""),
		Port:                    getEnv("PORT", "3000"),
		SentryDSN:               getEnv("SENTRY_DSN", ""),
		Environment:             getEnv("APP_ENV", "development"),
		DatabaseDSN:             dbDSN,
		RedisURL:                getEnv("REDIS_URL", "redis://localhost:6379/0"),
		RevenueCatWebhookSecret: rcSecret,
		AWSRegion:               getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:          getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:      getEnv("AWS_SECRET_ACCESS_KEY", ""),
		JWTSecret:               jwtSecret,
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
