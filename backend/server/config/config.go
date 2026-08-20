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
	OpenAIAPIKey            string
	FCMKey                  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		R2AccountID:             getEnv("R2_ACCOUNT_ID", "default_account_id"),
		R2AccessKeyID:           getEnv("R2_ACCESS_KEY_ID", "default_access_key"),
		R2SecretAccessKey:       getEnv("R2_SECRET_ACCESS_KEY", "default_secret_key"),
		R2BucketName:            getEnv("R2_BUCKET_NAME", "qlove-bucket"),
		Port:                    getEnv("PORT", "3000"),
		SentryDSN:               getEnv("SENTRY_DSN", ""),
		Environment:             getEnv("APP_ENV", "development"),
		DatabaseDSN:             getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=qlove port=5432 sslmode=disable"),
		RedisURL:                getEnv("REDIS_URL", "redis://localhost:6379/0"),
		RevenueCatWebhookSecret: getEnv("REVENUECAT_WEBHOOK_SECRET", "secret123"),
		AWSRegion:               getEnv("AWS_REGION", "ap-southeast-1"),
		AWSAccessKeyID:          getEnv("AWS_ACCESS_KEY_ID", "default_aws_access_key"),
		AWSSecretAccessKey:      getEnv("AWS_SECRET_KEY", "default_aws_secret_key"),
		JWTSecret:               jwtSecret,
		OpenAIAPIKey:            getEnv("OPENAI_API_KEY", ""),
		FCMKey:                  getEnv("FCM_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
